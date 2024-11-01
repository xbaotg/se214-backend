package router

import (
	authController "be/api/controller/auth"
	userController "be/api/controller/user"
	"be/api/middleware"
	"be/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine, app *bootstrap.App) {
	// public routes
	publicRouter := r.Group("")

	// auth routes
	publicRouter.POST("/register", func(c *gin.Context) {
		authController.Register(c, app)
	})
	publicRouter.POST("/login", func(c *gin.Context) {
		authController.Login(c, app)
	})

	// ----------------

	// protected routes
	protectedRouter := r.Group("")
	protectedRouter.Use(func(ctx *gin.Context) {
		middleware.SessionMiddleware(app)(ctx)
	})

	// auth routes
	protectedRouter.POST("/logout", func(c *gin.Context) {
		authController.Logout(c, app)
	})
	protectedRouter.POST("/refresh-token", func(c *gin.Context) {
		authController.RefreshToken(c, app)
	})

	// user routes
	protectedRouter.POST("/user/info", func(c *gin.Context) {
		userController.GetUserInfo(c, app)
	})
	protectedRouter.POST("/user/change-password", func(c *gin.Context) {
		userController.ChangePassUser(c, app)
	})
}
