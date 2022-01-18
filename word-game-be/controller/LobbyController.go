package controller

import (
	"context"
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/constant"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably-labs/word-game/word-game-be/util"
	"github.com/ably/ably-go/ably"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"time"
)

type LobbyController struct {
	db   *gorm.DB
	ably *ably.Realtime
}

func NewLobbyController(e *echo.Group, db *gorm.DB, ably *ably.Realtime) *LobbyController {

	lc := LobbyController{
		db:   db,
		ably: ably,
	}

	e.GET("", lc.GetLobbies)
	e.POST("", lc.PostLobby, middleware.RequireAuthorisation)

	// Endpoints which require a valid lobby
	lobbyGroup := e.Group("/:id", middleware.RequireAuthorisation, middleware.ValidateLobby)
	lobbyGroup.GET("", lc.GetLobby)
	lobbyGroup.GET("/thumbnail", lc.GetLobbyThumbnail)
	lobbyGroup.PUT("/member", lc.PutMember)
	lobbyGroup.GET("/member", lc.GetMembers, middleware.RequireLobbyMember)
	lobbyGroup.DELETE("/member", lc.DeleteMember, middleware.RequireLobbyMember)
	lobbyGroup.PATCH("/member", lc.PatchMember, middleware.RequireLobbyMember)

	return &lc
}

func (lc *LobbyController) GetLobbies(c echo.Context) error {
	var lobbies []model.Lobby
	lc.db.Preload("GameType").Preload("Creator").Find(&lobbies)
	return c.JSON(200, lobbies)
}

func (lc *LobbyController) GetLobby(c echo.Context) error {
	// TODO: Private lobbies should require membership
	return c.JSON(200, c.Get("lobby"))
}

func (lc *LobbyController) PostLobby(c echo.Context) error {
	createLobby := entity.CreateLobby{}
	err := c.Bind(&createLobby)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}
	user := c.Get("user").(*model.User)
	tileBag := util.NewTileBag()
	ownerDeck := util.TakeFromBag(7, &tileBag)
	newLobby := &model.Lobby{
		Name:           createLobby.Name,
		CreatorID:      user.ID,
		CreatedAt:      time.Now(),
		State:          model.LobbyWaitingForPlayers,
		Private:        createLobby.Private,
		Joinable:       true,
		CurrentPlayers: 1,
		MaxPlayers:     4,
		GameTypeID:     1,
		PlayerTurnID:   user.ID,
		Board:          util.NewBoard(),
		Bag:            tileBag,
	}

	err = lc.db.Create(newLobby).Error

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}

	newLobbyMember := &model.LobbyMember{
		UserID:     *user.ID,
		LobbyID:    *newLobby.ID,
		MemberType: constant.MemberTypePlayer,
		Deck: entity.SquareSet{
			Squares: &ownerDeck,
			Width:   9,
			Height:  1,
		},
	}

	err = lc.db.Create(newLobbyMember).Error

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}

	// Send a notification of a new lobby for public lobbies
	if !newLobby.Private {
		_ = lc.ably.Channels.Get("lobby-list").Publish(context.Background(), "lobbyAdd", newLobby)
	}

	return c.JSON(201, newLobby)
}

func (lc *LobbyController) GetLobbyThumbnail(c echo.Context) error {
	// TODO: Store a temporary thumbnail somewhere here that is occasionally regenerated
	return c.File("/home/peter/IdeaProjects/word-game/word-game-be/static/thumbnail-placeholder.png")
}

func (lc *LobbyController) PutMember(c echo.Context) error {
	putUser := entity.PutMember{}
	err := c.Bind(&putUser)
	if err != nil || (putUser.Type != constant.MemberTypePlayer && putUser.Type != constant.MemberTypeSpectator) {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobby := c.Get("lobby").(*model.Lobby)
	user := c.Get("user").(*model.User)

	// TODO: Codes
	if lobby.Private {
		return c.JSON(404, entity.ErrLobbyNotFound)
	}

	if putUser.Type == constant.MemberTypePlayer {
		if !lobby.Joinable {
			return c.JSON(403, entity.ErrForbidden)
		}

		if lobby.CurrentPlayers >= lobby.MaxPlayers {
			return c.JSON(403, entity.ErrLobbyFull)
		}
	}

	lobbyMember := &model.LobbyMember{
		UserID:     *user.ID,
		LobbyID:    *lobby.ID,
		MemberType: putUser.Type,
	}

	if lobbyMember.MemberType == constant.MemberTypePlayer {
		newDeck := util.TakeFromBag(7, &lobby.Bag)
		lobbyMember.Deck = entity.SquareSet{
			Squares: &newDeck,
			Width:   9,
			Height:  1,
		}
	}

	err = lc.db.Save(lobbyMember).Error

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}

	if lobbyMember.MemberType == constant.MemberTypePlayer {
		lobby.CurrentPlayers++
		lc.db.Save(lobby)
	}

	_ = util.PublishLobbyMessage(lc.ably, lobby.ID, "message", entity.ChatSent{
		Message: fmt.Sprintf("%s joined the game", user.Name),
		Author:  "system",
	})
	_ = util.PublishLobbyMessage(lc.ably, lobby.ID, "memberAdd", lobbyMember)

	return c.NoContent(204)
}

func (lc *LobbyController) DeleteMember(c echo.Context) error {
	delMember := entity.DeleteMember{}
	err := c.Bind(&delMember)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobby := c.Get("lobby").(*model.Lobby)
	user := c.Get("user").(*model.User)
	// If not lobby creator, only allow removing yourself
	if lobby.CreatorID != user.ID || delMember.UserID != user.ID {
		return c.JSON(403, entity.ErrForbidden)
	}

	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)

	err = lc.db.Delete(lobbyMember).Error
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}

	if lobbyMember.MemberType == constant.MemberTypePlayer {
		lobby.CurrentPlayers--
		lc.db.Save(lobby)
	}

	_ = util.PublishLobbyMessage(lc.ably, lobby.ID, "message", entity.ChatSent{
		Message: fmt.Sprintf("%s left the game", user.Name), // TODO: This shows the remover's name if another user was kicked
		Author:  "system",
	})
	_ = util.PublishLobbyMessage(lc.ably, lobby.ID, "memberRemove", lobbyMember)

	return c.NoContent(204)
}

func (lc *LobbyController) PatchMember(c echo.Context) error {
	putMember := entity.PutMember{}
	err := c.Bind(&putMember)
	if err != nil || (putMember.Type != constant.MemberTypePlayer && putMember.Type != constant.MemberTypeSpectator) {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobbyMember := c.Get("lobbyMember").(*model.LobbyMember)
	lobbyMember.MemberType = putMember.Type
	err = lc.db.Updates(&lobbyMember).Error

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, err)
	}

	// TODO Updating a member to a player doesn't change the player count

	_ = util.PublishLobbyMessage(lc.ably, &lobbyMember.LobbyID, "memberUpdate", lobbyMember)

	return c.NoContent(204)
}

func (lc *LobbyController) GetMembers(c echo.Context) error {
	var members []model.LobbyMember
	lobby := c.Get("lobby").(*model.Lobby)
	err := lc.db.Preload("User").Where("lobby_id = ?", lobby.ID).Find(&members).Error
	if err != nil {
		return c.JSON(500, entity.ErrDatabaseError)
	}

	return c.JSON(200, members)
}
