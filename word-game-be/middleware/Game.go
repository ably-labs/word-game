package middleware

import (
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/labstack/echo/v4"
)

func RequireTurn(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lobby := c.Get("lobby").(*model.Lobby)
		user := c.Get("user").(*model.User)
		if lobby.PlayerTurnID != user.ID {
			return c.JSON(403, entity.ErrNotYourTurn)
		}
		return handlerFunc(c)
	}
}
