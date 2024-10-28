package main

import (
	route "be/api/router"
	"be/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.NewApp()

	// run the server
	r := gin.Default()
	route.SetupRoute(r, app)
	r.Run(app.Config.ServerAddr)
}
