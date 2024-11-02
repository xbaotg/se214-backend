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
	r.Run(app.Config.ServerAddr)
}
