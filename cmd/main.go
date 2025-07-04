package main

import (
	"net/http"

	"github.com/appu900/deerDB/auth"
	"github.com/appu900/deerDB/db"
	"github.com/appu900/deerDB/types"
	"github.com/labstack/echo/v4"
)

func main() {
	config := &types.Config{
		Port:      "8080",
		JWTSecret: "my-secret",
	}

	authservice, err := auth.NewAuthService(config)
	if err != nil {
		panic(err)
	}
	db.Testdatabase()
	defer authservice.Close()
	e := echo.New()
	authHandler := auth.NewAuthHttpHandler(authservice)

	e.GET("/ping", Ping)
	e.POST("/register", authHandler.HandleUserRegistration)
	e.POST("/login", authHandler.HandleUserLogin)
	e.GET("/user/profile", authHandler.HandlerUserProfile, auth.RequireAuth())
	e.Logger.Fatal(e.Start(":" + config.Port))
}

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "its working")
}
