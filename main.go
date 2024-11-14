package main

import (
	"be/api/router"
	"be/bootstrap"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.NewApp()

	// run the server
	r := gin.Default()

	// cors
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	router.SetupRoute(r, app)
	router.SetupSwaggerRoute(r)

	err := r.Run(app.Config.ServerAddr)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
