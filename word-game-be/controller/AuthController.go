package controller

import (
	"encoding/gob"
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"os"
)

type AuthController struct {
	Auth *webauthn.WebAuthn
	db   *gorm.DB
}

func NewAuthController(e *echo.Group, db *gorm.DB) *AuthController {
	feRoot := os.Getenv("FRONTEND_BASE_URL")
	//beRoot := os.Getenv("BACKEND_BASE_URL")
	web, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Word Game",
		RPID:          "localhost",
		RPOrigin:      fmt.Sprintf("http://%s/", feRoot),
		RPIcon:        fmt.Sprintf("http://%s/letter_w.png", feRoot),
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			RequireResidentKey: protocol.ResidentKeyUnrequired(),
			UserVerification:   protocol.VerificationDiscouraged,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	gob.Register(webauthn.SessionData{})
	gob.Register(model.User{})

	ac := AuthController{
		Auth: web,
		db:   db,
	}

	e.POST("register/start", ac.PostStartRegister)
	e.POST("register/confirm", ac.PostCompleteRegister)

	return &ac
}

func (ac *AuthController) PostStartRegister(c echo.Context) error {
	body := entity.Register{}
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(400, entity.Error{Err: "Invalid Input"})
	}
	uid, _ := uuid.NewUUID()
	sess, _ := session.Get("session", c)
	userId := uid.ID()
	newUser := &model.User{
		ID:   &userId,
		Name: body.Nickname,
	}
	options, sessionData, _ := ac.Auth.BeginRegistration(newUser)
	sess.Values["register_session"] = sessionData
	sess.Values["register_user"] = *newUser
	err = sess.Save(c.Request(), c.Response())
	fmt.Println(err)
	return c.JSON(200, options)
}

func (ac *AuthController) PostCompleteRegister(c echo.Context) error {
	body, err := protocol.ParseCredentialCreationResponseBody(c.Request().Body)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(400, err)
	}
	sess, _ := session.Get("session", c)
	newUser, ok := sess.Values["register_user"].(model.User)
	if !ok {
		return c.JSON(400, entity.Error{Err: "Bad session"})
	}
	fmt.Println(newUser)
	credential, err := ac.Auth.CreateCredential(&newUser, sess.Values["register_session"].(webauthn.SessionData), body)
	if err != nil {
		return c.JSON(400, err)
	}
	fmt.Println(credential)
	// TODO: Ably token
	newUser.Credentials = []webauthn.Credential{*credential}
	ac.db.Create(&newUser)
	return c.JSON(200, newUser)
}
