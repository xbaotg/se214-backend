package router

import (
	authController "be/api/controller/auth"
	coursesController "be/api/controller/courses"
	departmentsController "be/api/controller/departments"
	usersController "be/api/controller/users"
	"be/api/middleware"
	"be/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine, app *bootstrap.App) {
	// public routes
	publicRouter := r.Group("")
	publicRouter.Use(func(ctx *gin.Context) {
		middleware.SessionMiddleware(app, false)(ctx)
	})

	// auth routes
	publicRouter.POST("/register", func(c *gin.Context) {
		authController.Register(c, app)
	})
	publicRouter.POST("/login", func(c *gin.Context) {
		authController.Login(c, app)
	})

	// department routes
	publicRouter.GET("/department/list", func(c *gin.Context) {
		departmentsController.ListDepartment(c, app)
	})

	// ----------------

	// protected routes
	protectedRouter := r.Group("")
	protectedRouter.Use(func(ctx *gin.Context) {
		middleware.SessionMiddleware(app, true)(ctx)
	})

	// auth routes
	protectedRouter.GET("/logout", func(c *gin.Context) {
		authController.Logout(c, app)
	})
	protectedRouter.POST("/refresh-token", func(c *gin.Context) {
		authController.RefreshToken(c, app)
	})

	// user routes
	protectedRouter.GET("/user/info", func(c *gin.Context) {
		usersController.GetUserInfo(c, app)
	})
	protectedRouter.PATCH("/user/change-password", func(c *gin.Context) {
		usersController.ChangePassUser(c, app)
	})

	// course routes
	protectedRouter.POST("/course/create", func(c *gin.Context) {
		coursesController.CreateCourse(c, app)
	})

	// department routes
	protectedRouter.POST("/department/create", func(c *gin.Context) {
		departmentsController.CreateDepartment(c, app)
	})
}
