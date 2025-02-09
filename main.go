package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sj0n/echo-api/pkg/routes"
)

func main() {
	server := echo.New()

	server.Use(middleware.Logger())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://popdex.harizdan.xyz", "popdex.pages.dev"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "cache-control", "if-none-match"},
	}))

	server.GET("/pokemon/:name", routes.GetPokemonProfile)
	server.GET("/pokemon/:name/moves", routes.GetPokemonMoves)
	server.GET("/pokemon/:name/locations", routes.GetPokemonLocations)

	server.Logger.Fatal(server.Start(":8081"))
}
