package main

import (
	"be/api/router"
	"be/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.NewApp()

	// run the server
	r := gin.Default()
	router.SetupRoute(r, app)
	router.SetupSwaggerRoute(r)

	err := r.Run(app.Config.ServerAddr)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
