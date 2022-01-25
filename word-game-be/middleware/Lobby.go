package middleware

import (
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strconv"
)

func ValidateLobby(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := c.Get("db").(*gorm.DB)
		lobbyId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(400, entity.ErrInvalidLobby)
		}

		castLobbyId := int64(lobbyId)

		lobby := model.Lobby{ID: &castLobbyId}
		err = db.Preload("GameType").Find(&lobby).Error

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

func RequireLobbyOwner(handleFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lobby := c.Get("lobby").(*model.Lobby)
		user := c.Get("user").(*model.User)
		if *lobby.CreatorID != *user.ID {

			return c.JSON(403, entity.ErrForbidden)
		}
		return handleFunc(c)
	}
}

func RequireLobbyMember(handleFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lobby := c.Get("lobby").(*model.Lobby)
		user := c.Get("user").(*model.User)
		db := c.Get("db").(*gorm.DB)
		lobbyMember := model.LobbyMember{
			UserID:  *user.ID,
			LobbyID: *lobby.ID,
		}

		err := db.First(&lobbyMember).Error

		if err == gorm.ErrRecordNotFound {
			return c.JSON(403, entity.ErrForbidden)
		} else if err != nil {
			return c.JSON(500, entity.ErrDatabaseError)
		}

		c.Set("lobbyMember", &lobbyMember)

		return handleFunc(c)
	}
}
