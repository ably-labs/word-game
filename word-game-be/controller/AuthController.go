package controller

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/entity"
	"github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably/ably-go/ably"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

type AuthController struct {
	Auth *webauthn.WebAuthn
	db   *gorm.DB
	ably *ably.Realtime
}

func NewAuthController(e *echo.Group, db *gorm.DB, ably *ably.Realtime) *AuthController {
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
		ably: ably,
	}

	// Registration
	e.POST("/register/start", ac.PostStartRegister, middleware.DisallowAuthorisation)
	e.POST("/register/confirm", ac.PostCompleteRegister, middleware.DisallowAuthorisation)

	// Login
	e.POST("/login/start", ac.PostStartLogin, middleware.DisallowAuthorisation)
	e.POST("/login/confirm", ac.PostCompleteLogin, middleware.DisallowAuthorisation)

	e.GET("/token", ac.GetAblyToken, middleware.RequireAuthorisation)
	e.GET("/me", ac.GetMe, middleware.RequireAuthorisation)
	return &ac
}

// PostStartRegister initiates the WebAuthn request on the client
func (ac *AuthController) PostStartRegister(c echo.Context) error {
	body := entity.Register{}
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}
	// Create the new user model early to check the nickname is unique
	newUser := &model.User{
		Name: body.Nickname,
	}

	if newUser.Exists(ac.db) {
		return c.JSON(409, entity.ErrUserExists)
	}

	// Generate a new UUID and store it as the user's ID
	uid := uuid.New()
	userId := uid.ID()
	newUser.ID = &userId

	options, sessionData, err := ac.Auth.BeginRegistration(newUser)

	if err != nil {
		return c.JSON(400, err)
	}

	sess, _ := session.Get("session", c)
	sess.Values["register_session"] = sessionData
	sess.Values["register_user"] = *newUser
	err = sess.Save(c.Request(), c.Response())

	return c.JSON(200, options)
}

// PostCompleteRegister completes the registration and creates the user
func (ac *AuthController) PostCompleteRegister(c echo.Context) error {
	// Parse the incoming request as a CredentialCreationResponseBody
	body, err := protocol.ParseCredentialCreationResponseBody(c.Request().Body)
	if err != nil {
		return c.JSON(400, err)
	}

	// Retrieve the session and registration values
	sess, _ := session.Get("session", c)
	newUser, userOk := sess.Values["register_user"].(model.User)
	registerSession, sessionOk := sess.Values["register_session"].(webauthn.SessionData)
	// If the values aren't there, the user hasn't initiated registration
	if !userOk || !sessionOk {
		return c.JSON(400, entity.ErrBadSession)
	}
	// Create the credential
	credential, err := ac.Auth.CreateCredential(&newUser, registerSession, body)
	if err != nil {
		return c.JSON(400, err)
	}
	credJson, _ := json.Marshal([]webauthn.Credential{*credential})
	newUser.Credentials = credJson

	err = ac.db.Create(&newUser).Error
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrUserCreateFailed)
	}

	// Clear the session data
	sess.Values["register_session"] = nil
	sess.Values["register_user"] = nil
	sess.Values["user_id"] = newUser.ID
	_ = sess.Save(c.Request(), c.Response())

	return c.JSON(200, newUser)
}

func (ac *AuthController) PostStartLogin(c echo.Context) error {
	body := entity.Register{}
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(400, entity.ErrInvalidInput)
	}

	user := model.User{
		Name: body.Nickname,
	}

	err = ac.db.Where(&user).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(404, entity.ErrUserDoesNotExist)
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(500, entity.ErrDatabaseError)
	}
	fmt.Println(user)

	_ = json.Unmarshal(user.Credentials, &user.CredentialsObj)

	options, sessionData, err := ac.Auth.BeginLogin(&user)

	if err != nil {
		return c.JSON(400, err)
	}

	sess, _ := session.Get("session", c)
	sess.Values["login_session"] = sessionData
	sess.Values["login_user"] = user
	err = sess.Save(c.Request(), c.Response())

	return c.JSON(200, options)
}

func (ac *AuthController) PostCompleteLogin(c echo.Context) error {
	body, err := protocol.ParseCredentialRequestResponseBody(c.Request().Body)
	if err != nil {
		return c.JSON(400, err)
	}

	sess, _ := session.Get("session", c)
	sessionData, sessOk := sess.Values["login_session"].(webauthn.SessionData)
	user, userOk := sess.Values["login_user"].(model.User)

	if !sessOk || !userOk {
		return c.JSON(400, entity.ErrDatabaseError)
	}

	sess.Values["login_session"] = nil
	sess.Values["login_user"] = nil

	_, err = ac.Auth.ValidateLogin(&user, sessionData, body)

	if err != nil {
		_ = sess.Save(c.Request(), c.Response())
		return c.JSON(400, err)
	}

	user.LastActive = time.Now()
	ac.db.Save(user)

	sess.Values["user_id"] = user.ID
	_ = sess.Save(c.Request(), c.Response())

	return c.JSON(200, user)
}

func (ac *AuthController) GetAblyToken(c echo.Context) error {
	user := c.Get("user").(*model.User)

	req, err := ac.ably.Auth.CreateTokenRequest(&ably.TokenParams{
		ClientID: strconv.Itoa(int(*user.ID)),
	})
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, req)
}

func (ac *AuthController) GetMe(c echo.Context) error {
	return c.JSON(200, c.Get("user"))
}
