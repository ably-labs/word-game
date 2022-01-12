package controller

import (
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LobbyController struct {
	db *gorm.DB
}

func NewLobbyController(e *echo.Group, db *gorm.DB) *LobbyController {

	lc := LobbyController{
		db: db,
	}

	e.GET("/", lc.GetLobbies)

	return &lc
}

func (lc *LobbyController) GetLobbies(c echo.Context) error {
	var lobbies []model.Lobby
	lc.db.Find(&lobbies)
	return c.JSON(200, lobbies)
}
