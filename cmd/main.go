package main

import (
	"fmt"
	"game-server-example/cmd/di"
	"game-server-example/cmd/di/provider"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func start() error {
	var gameServer *provider.GameServer
	var cleanup func()
	var err error

	gameServer, cleanup, err = di.Inject()
	if err != nil {
		return err
	}
	defer cleanup()

	fmt.Println("ENABLE_CACHE_REPOSITORY:", os.Getenv("ENABLE_CACHE_REPOSITORY"))

	e := echo.New()
	gameServer.Router.RegisterRoute(e)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	return nil
}
