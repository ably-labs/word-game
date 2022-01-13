package middleware

import (
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AuthoriseUser(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := c.Get("db").(*gorm.DB)
		sess, _ := session.Get("session", c)
		userId, userOk := sess.Values["user_id"].(uint32)

		if !userOk {
			// Unauthed is allowed by default
			return handlerFunc(c)

		}

		user := model.User{
			ID: &userId,
		}

		err := db.First(&user).Error

		if err == gorm.ErrRecordNotFound {
			// User doesn't exist anymore
			return handlerFunc(c)
		} else if err != nil {
			return c.JSON(500, entity.ErrDatabaseError)
		}

		c.Set("user", &user)

		return handlerFunc(c)
	}
}

func RequireAuthorisation(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") == nil {
			return c.JSON(401, entity.ErrUnauthorised)
		}
		return handlerFunc(c)
	}
}

func DisallowAuthorisation(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") != nil {
			return c.JSON(403, entity.ErrLoggedIn)
		}
		return handlerFunc(c)
	}
}
