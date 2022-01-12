package main

import (
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/controller"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	requiredEnvVariables("DB_DSN", "SESSION_SECRET", "FRONTEND_BASE_URL", "BACKEND_BASE_URL")

	// Open the DB and migrate the models
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(model.User{})

	if err != nil {
		log.Fatalln(err)
	}

	// Create the web server and routes
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// TODO: This should be actual origins
			return true, nil
		},
		AllowCredentials: true,
	}))
	// Initialise the session storage.
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	controller.NewAuthController(e.Group("auth/"), db)

	// Start the web server
	e.Logger.Fatal(e.Start(":3001"))
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
