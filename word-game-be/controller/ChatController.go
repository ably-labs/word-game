package controller

import (
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably-labs/word-game/word-game-be/util"
	"github.com/ably/ably-go/ably"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"time"
)

type ChatController struct {
	db   *gorm.DB
	ably *ably.Realtime
}

func NewChatController(e *echo.Group, db *gorm.DB, client *ably.Realtime) *ChatController {

	cc := &ChatController{
		db:   db,
		ably: client,
	}

	e.Use(middleware.RequireAuthorisation)

	lobbyGroup := e.Group("/:id", middleware.ValidateLobby, middleware.RequireLobbyMember)
	lobbyGroup.GET("", cc.GetChatHistory)
	lobbyGroup.POST("", cc.PostChatMessage)

	return cc
}

func (cc *ChatController) PostChatMessage(c echo.Context) error {
	chatInput := entity.SendChat{}
	err := c.Bind(&chatInput)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	lobby := c.Get("lobby").(*model.Lobby)
	user := c.Get("user").(*model.User)

	message := model.Message{
		AuthorID:  user.ID,
		LobbyID:   lobby.ID,
		Message:   chatInput.Message,
		Timestamp: time.Now(),
	}

	err = cc.db.Create(&message).Error
	if err != nil {
		return c.JSON(500, entity.ErrDatabaseError)
	}

	err = util.PublishLobbyMessage(cc.ably, lobby.ID, "message", entity.ChatSent{
		Message: chatInput.Message,
		Author:  user.Name,
	})

	return c.NoContent(204)
}

func (cc *ChatController) GetChatHistory(c echo.Context) error {
	lobby := c.Get("lobby").(*model.Lobby)
	var messages []model.Message

	err := cc.db.Where("lobby_id = ?", lobby.ID).Preload("Author").Limit(50).Order("Timestamp").Find(&messages).Error

	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}

	cleanedMessages := make([]entity.ChatSent, len(messages))
	for i, message := range messages {
		cleanedMessages[i] = entity.ChatSent{
			Message: message.Message,
			Author:  message.Author.Name,
		}
	}

	return c.JSON(200, cleanedMessages)
}
