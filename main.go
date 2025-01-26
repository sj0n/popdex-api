package main

import (
	"github.com/labstack/echo/v4"
	
	"github.com/sj0n/echo-api/pkg/routes"
)

func main() {
	server := echo.New()

	server.GET("/pokemon/:name", routes.GetPokemonProfile)
	server.GET("/pokemon/:name/moves", routes.GetPokemonMoves)
	server.GET("/pokemon/:name/locations", routes.GetPokemonLocations)

	server.Logger.Fatal(server.Start(":8081"))
}
