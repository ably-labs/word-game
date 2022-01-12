package controller

import (
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strconv"
)

type LobbyController struct {
	db *gorm.DB
}

func NewLobbyController(e *echo.Group, db *gorm.DB) *LobbyController {

	lc := LobbyController{
		db: db,
	}

	e.GET("", lc.GetLobbies)

	// Endpoints which require a valid lobby
	lobbyGroup := e.Group("/:id", middleware.RequireAuthorisation, lc.MwValidateLobby)
	lobbyGroup.GET("/thumbnail", lc.GetLobbyThumbnail)
	lobbyGroup.PUT("/player", lc.PutPlayer)
	lobbyGroup.PUT("/spectator", lc.PutSpectator)

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
			return c.JSON(400, entity.Error{Err: "Lobby ID must be an integer"})
		}

		castLobbyId := uint32(lobbyId)

		lobby := model.Lobby{ID: &castLobbyId}
		err = lc.db.Find(&lobby).Error

		if err == gorm.ErrRecordNotFound {
			return c.JSON(404, entity.Error{Err: "Lobby Not Found"})
		}
		if err != nil {
			return c.JSON(500, entity.Error{Err: "Database Error"})
		}
		c.Set("lobby", &lobby)

		return handlerFunc(c)
	}
}

func (lc *LobbyController) MwRequireLobbyOwner(handleFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lobby := c.Get("lobby").(*model.Lobby)
		user := c.Get("user").(*model.User)
		//if lobby.CreatorID !=
		return nil
	}
}

func (lc *LobbyController) GetLobbyThumbnail(c echo.Context) error {
	// TODO: Store a temporary thumbnail somewhere here that is occasionally regenerated
	return c.File("/home/peter/IdeaProjects/word-game/word-game-be/static/thumbnail-placeholder.png")
}

func (lc *LobbyController) PutPlayer(c echo.Context) error {
	return nil
}

func (lc *LobbyController) PutSpectator(c echo.Context) error {
	return nil
}
