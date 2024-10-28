package router

import (
	"be/api/controller"
	"be/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine, app *bootstrap.App) {
	// auth routes
	r.POST("/register", func(c *gin.Context) {
		controller.Register(c, app)
	})
	r.POST("/login", func(c *gin.Context) {
		controller.Login(c, app)
	})
	r.POST("/logout", func(c *gin.Context) {
		controller.Logout(c, app)
	})
	r.POST("/refresh-token", func(c *gin.Context) {
		controller.RefreshToken(c, app)
	})

	// user routes
}
