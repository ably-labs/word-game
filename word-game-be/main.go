package main

import (
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/controller"
	localMiddleware "github.com/ably-labs/word-game/word-game-be/middleware"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably/ably-go/ably"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	requiredEnvVariables("DB_DSN", "SESSION_SECRET", "FRONTEND_BASE_URL", "BACKEND_BASE_URL", "ABLY_KEY")

	// Open the DB and migrate the models
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(model.User{}, model.GameType{}, model.Lobby{}, model.LobbyMember{}, model.Message{})

	if err != nil {
		log.Fatalln(err)
	}

	// Save the standard type
	standardType := model.GameType{
		ID:               1,
		Name:             "Standard",
		TileBagSize:      100,
		PlayerDeckSize:   9,
		PlayerTileCount:  7,
		EnableBlankTiles: true,
		BoardWidth:       15,
		BoardHeight:      15,
		BonusTilePattern: model.Regular,
		StartAnywhere:    false,
	}
	err = db.Save(&standardType).Error
	if err != nil {
		log.Fatalln(err)
	}

	ablyClient := initAbly()

	// Create the web server and routes
	e := echo.New()

	// Inject the DB instance into echo to use with middlewares
	e.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return handlerFunc(c)
		}
	})

	// Setup CORS Headers
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// TODO: This should be actual origins
			return true, nil
		},
		AllowCredentials: true,
	}))

	// Initialise the session storage.
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	e.Use(localMiddleware.AuthoriseUser)

	controller.NewAuthController(e.Group("auth"), db, ablyClient)
	controller.NewLobbyController(e.Group("lobby"), db, ablyClient)
	controller.NewChatController(e.Group("chat"), db, ablyClient)
	controller.NewGameController(e.Group("game"), db, ablyClient)

	// Start the web server
	e.Logger.Fatal(e.Start(":3001"))
}

func initAbly() *ably.Realtime {

	client, err := ably.NewRealtime(ably.WithKey(os.Getenv("ABLY_KEY")))
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

// Small util just to check every environment variable exists that is required
func requiredEnvVariables(variables ...string) {
	isMissing := false
	for _, variable := range variables {
		if os.Getenv(variable) == "" {
			isMissing = true
			fmt.Printf("Missing required environment variable %s\n", variable)
		}
	}

	if isMissing {
		fmt.Println("One or more required environment variables is missing. Please set them and restart.")
		os.Exit(1)
	}
}
