package controller

import (
	"context"
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/constant"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably/ably-go/ably"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strconv"
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

	// Endpoints which require a valid lobby
	lobbyGroup := e.Group("/:id", middleware.RequireAuthorisation, lc.MwValidateLobby)
	lobbyGroup.GET("/thumbnail", lc.GetLobbyThumbnail)
	lobbyGroup.PUT("/member", lc.PutMember)
	lobbyGroup.DELETE("/member", lc.DeleteMember)
	lobbyGroup.PATCH("/member", lc.PatchMember)

	return &lc
}

func (lc *LobbyController) GetLobbies(c echo.Context) error {
	var lobbies []model.Lobby
	lc.db.Preload("GameType").Find(&lobbies)
	return c.JSON(200, lobbies)
}

func (lc *LobbyController) MwValidateLobby(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lobbyId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(400, entity.ErrInvalidLobby)
		}

		castLobbyId := uint32(lobbyId)

		lobby := model.Lobby{ID: &castLobbyId}
		err = lc.db.Find(&lobby).Error

		if err == gorm.ErrRecordNotFound {
			return c.JSON(404, entity.ErrLobbyNotFound)
		}
		if err != nil {
			return c.JSON(500, entity.ErrDatabaseError)
		}
		c.Set("lobby", &lobby)

		return handlerFunc(c)
	}
}

func (lc *LobbyController) MwRequireLobbyOwner(handleFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lobby := c.Get("lobby").(*model.Lobby)
		user := c.Get("user").(*model.User)
		if lobby.CreatorID != user.ID {
			return c.JSON(403, entity.ErrForbidden)
		}
		return handleFunc(c)
	}
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

	if !lobby.Joinable && putUser.Type == constant.MemberTypePlayer {
		return c.JSON(403, entity.ErrForbidden)
	}

	lobbyMember := &model.LobbyMember{
		UserID:     *user.ID,
		LobbyID:    *lobby.ID,
		MemberType: putUser.Type,
	}

	err = lc.db.Save(lobbyMember).Error

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}

	_ = lc.publishLobbyMessage(lobby.ID, "memberAdd", lobbyMember)

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

	lobbyMember := &model.LobbyMember{
		UserID:  *user.ID,
		LobbyID: *lobby.ID,
	}

	err = lc.db.Delete(lobbyMember).Error
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}

	_ = lc.publishLobbyMessage(lobby.ID, "memberRemove", lobbyMember)

	return c.NoContent(204)
}

func (lc *LobbyController) PatchMember(c echo.Context) error {
	putMember := entity.PutMember{}
	err := c.Bind(&putMember)
	if err != nil || (putMember.Type != constant.MemberTypePlayer && putMember.Type != constant.MemberTypeSpectator) {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobby := c.Get("lobby").(*model.Lobby)
	user := c.Get("user").(*model.User)
	lobbyMember := model.LobbyMember{
		UserID:     *user.ID,
		LobbyID:    *lobby.ID,
		MemberType: putMember.Type,
	}

	err = lc.db.Updates(&lobbyMember).Error

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, err)
	}

	_ = lc.publishLobbyMessage(lobby.ID, "memberUpdate", lobbyMember)

	return c.NoContent(204)
}

func (lc *LobbyController) publishLobbyMessage(lobby *uint32, name string, message interface{}) error {
	return lc.ably.Channels.Get(fmt.Sprintf("lobby-%d", lobby)).Publish(context.Background(), name, message)
}
