package router

import (
	"be/api/controller"
	"be/api/middleware"
	"be/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine, app *bootstrap.App) {
	// public routes
	publicRouter := r.Group("")

	// auth routes
	publicRouter.POST("/register", func(c *gin.Context) {
		controller.Register(c, app)
	})
	publicRouter.POST("/login", func(c *gin.Context) {
		controller.Login(c, app)
	})

	// ----------------

	// protected routes
	protectedRouter := r.Group("")
	protectedRouter.Use(func(ctx *gin.Context) {
		middleware.SessionMiddleware(app)(ctx)
	})

	// auth routes
	protectedRouter.POST("/logout", func(c *gin.Context) {
		controller.Logout(c, app)
	})
	protectedRouter.POST("/refresh-token", func(c *gin.Context) {
		controller.RefreshToken(c, app)
	})

	// user routes
	protectedRouter.POST("/user-info", func(c *gin.Context) {
		controller.GetUserInfo(c, app)
	})
}
