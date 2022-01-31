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
	"strconv"
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
	e.GET("/joined", lc.GetJoinedLobbies, middleware.RequireAuthorisation)

	e.GET("/types", lc.GetGameTypes)
	e.POST("/types", lc.PostGameType, middleware.RequireAuthorisation)

	// Endpoints which require a valid lobby
	lobbyGroup := e.Group("/:id", middleware.RequireAuthorisation, middleware.ValidateLobby)
	lobbyGroup.GET("", lc.GetLobby)
	lobbyGroup.DELETE("", lc.DeleteLobby, middleware.RequireLobbyOwner)
	lobbyGroup.PATCH("", lc.PatchLobby, middleware.RequireLobbyOwner)
	lobbyGroup.GET("/thumbnail", lc.GetLobbyThumbnail)
	lobbyGroup.PUT("/member", lc.PutMember)
	lobbyGroup.GET("/member", lc.GetMembers, middleware.RequireLobbyMember)
	lobbyGroup.DELETE("/member", lc.DeleteMember, middleware.RequireLobbyMember)
	lobbyGroup.PATCH("/member", lc.PatchMember, middleware.RequireLobbyMember)

	return &lc
}

func (lc *LobbyController) GetJoinedLobbies(c echo.Context) error {
	user := c.Get("user").(*model.User)

	lc.db.Where("user_id = ?", user.ID).Find(&user.LobbyMemberships)

	for i, lm := range *user.LobbyMemberships {
		(*user.LobbyMemberships)[i].LobbyIDStr = strconv.Itoa(int(lm.LobbyID))
	}

	return c.JSON(200, user.LobbyMemberships)
}

func (lc *LobbyController) PatchLobby(c echo.Context) error {
	lobbyUpdate := entity.UpdateLobby{}
	err := c.Bind(&lobbyUpdate)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	if lobbyUpdate.State != entity.LobbyInGame {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobby := c.Get("lobby").(*model.Lobby)

	lobby.State = lobbyUpdate.State

	err = lc.db.Save(lobby).Error

	if err != nil {
		return c.JSON(500, entity.ErrDatabaseError)
	}

	_ = util.PublishLobbyMessage(lc.ably, lobby.ID, "message", entity.ChatSent{
		Message: "Game has started!",
		Author:  "system",
	})

	_ = util.PublishLobbyMessage(lc.ably, lobby.ID, "lobbyUpdate", lobby)

	return c.JSON(200, lobby)

}

func (lc *LobbyController) DeleteLobby(c echo.Context) error {
	lobby := c.Get("lobby").(*model.Lobby)

	err := lc.db.Delete(&lobby).Error

	if err != nil {
		return c.JSON(500, err)
	}

	_ = util.PublishLobbyMessage(lc.ably, lobby.ID, "lobbyDeleted", nil)

	return c.NoContent(204)
}

func (lc *LobbyController) GetGameTypes(c echo.Context) error {
	var gameTypes []model.GameType
	lc.db.Find(&gameTypes)
	return c.JSON(200, gameTypes)
}

func (lc *LobbyController) PostGameType(c echo.Context) error {
	gameType := model.GameType{}
	err := c.Bind(&gameType)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	// Zero out the ID to prevent people setting their own
	gameType.ID = 0
	gameType.Visible = false

	err = lc.db.Save(&gameType).Error

	if err != nil {
		return c.JSON(400, entity.ErrDatabaseError)
	}

	return c.JSON(200, gameType)
}

func (lc *LobbyController) GetLobbies(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	var lobbies []model.Lobby
	query := lc.db.Preload("GameType").Preload("Creator").Order("state")

	if ok && user != nil {
		query = query.Where("private = false OR id IN (SELECT lobby_id FROM lobby_members WHERE user_id = ?)", *user.ID)
	} else {
		query = query.Where("private = false")
	}

	query.Find(&lobbies)
	for i := range lobbies {
		lobbies[i].IdStr = strconv.Itoa(int(*lobbies[i].ID))
	}

	return c.JSON(200, lobbies)
}

func (lc *LobbyController) GetLobby(c echo.Context) error {
	// TODO: Private lobbies should require membership
	lobby := c.Get("lobby").(*model.Lobby)
	user := c.Get("user").(*model.User)

	lobbyMember := model.LobbyMember{UserID: *user.ID, LobbyID: *lobby.ID}

	err := lc.db.Find(&lobbyMember).Error

	if err != nil {
		return c.JSON(500, entity.ErrDatabaseError)
	}

	if lobbyMember.MemberType == "" {
		return c.JSON(403, entity.ErrLobbyNotJoined)
	}

	return c.JSON(200, lobby)
}

func (lc *LobbyController) PostLobby(c echo.Context) error {
	createLobby := entity.CreateLobby{}
	err := c.Bind(&createLobby)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}
	user := c.Get("user").(*model.User)

	gameType := model.GameType{
		ID: createLobby.GameType,
	}

	err = lc.db.Find(&gameType).Error

	if err != nil {
		return c.JSON(400, entity.ErrDatabaseError)
	}

	tileBag := util.NewTileBag()
	ownerDeck := util.TakeFromBag(gameType.PlayerTileCount, &tileBag)
	newLobby := &model.Lobby{
		Name:           createLobby.Name,
		CreatorID:      user.ID,
		CreatedAt:      time.Now(),
		State:          entity.LobbyWaitingForPlayers,
		Private:        createLobby.Private,
		Joinable:       true,
		CurrentPlayers: 1,
		MaxPlayers:     4,
		GameTypeID:     createLobby.GameType,
		PlayerTurnID:   user.ID,
		Board:          util.NewBoard(gameType.BoardWidth, gameType.BoardHeight),
		Bag:            tileBag,
	}

	err = lc.db.Create(newLobby).Error

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}

	newLobby.IdStr = strconv.Itoa(int(*newLobby.ID))

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
	return c.File("/home/peter/wordgame/static/thumbnail-placeholder.png")
}

func (lc *LobbyController) PutMember(c echo.Context) error {
	putUser := entity.PutMember{}
	err := c.Bind(&putUser)
	if err != nil || (putUser.Type != constant.MemberTypePlayer && putUser.Type != constant.MemberTypeSpectator) {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobby := c.Get("lobby").(*model.Lobby)
	user := c.Get("user").(*model.User)

	// Check if the game is already in progress and the user is trying to join as a real user
	if lobby.State == entity.LobbyInGame && putUser.Type != constant.MemberTypeSpectator {
		return c.JSON(400, entity.ErrGameInProgress)
	}

	lobbyMember := model.LobbyMember{UserID: *user.ID, LobbyID: *lobby.ID}

	lc.db.Find(&lobbyMember)

	// Check if the user is already a player or if they are trying to join as the member they currently are
	if lobbyMember.MemberType == constant.MemberTypePlayer || lobbyMember.MemberType == putUser.Type {
		return c.JSON(400, entity.ErrLobbyJoined)
	}

	if putUser.Type == constant.MemberTypePlayer {
		if !lobby.Joinable {
			return c.JSON(403, entity.ErrForbidden)
		}

		if lobby.CurrentPlayers >= lobby.MaxPlayers {
			return c.JSON(403, entity.ErrLobbyFull)
		}
	}

	lobbyMember = model.LobbyMember{
		UserID:     *user.ID,
		LobbyID:    *lobby.ID,
		MemberType: putUser.Type,
	}

	// If they are a player, create a tile deck for them and take from the bag
	if lobbyMember.MemberType == constant.MemberTypePlayer {
		newDeck := util.TakeFromBag(7, &lobby.Bag)
		squares := make([]entity.Square, 9)
		lobbyMember.Deck = entity.SquareSet{
			Squares: &squares,
			Width:   9,
			Height:  1,
		}
		lobbyMember.Deck.AddTiles(newDeck)
	}

	err = lc.db.Save(&lobbyMember).Error

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

	lobbyMember.User = &model.DisplayUser{
		ID:   user.ID,
		Name: user.Name,
	}

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
