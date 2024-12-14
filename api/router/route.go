package router

import (
	authController "be/api/controller/auth"
	coursesController "be/api/controller/courses"
	departmentsController "be/api/controller/departments"
	usersController "be/api/controller/users"
	usersCourseController "be/api/controller/users/courses_managament"
	lecturerCourseController "be/api/controller/users/lecturer"
	globalController "be/api/controller/global"
	tuitionController "be/api/controller/tuition"
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
	publicRouter.POST("/register", func(c *gin.Context) { authController.Register(c, app) })
	publicRouter.POST("/login", func(c *gin.Context) { authController.Login(c, app) })

	// course routes
	publicRouter.GET("/course/list", func(c *gin.Context) { coursesController.ListCourses(c, app) })

	// department routes
	publicRouter.GET("/department/list", func(c *gin.Context) { departmentsController.ListDepartment(c, app) })

	// ----------------

	// protected routes
	protectedRouter := r.Group("")
	protectedRouter.Use(func(ctx *gin.Context) { middleware.SessionMiddleware(app, true)(ctx) })

	protectedRouter.POST("/global/state", func(c *gin.Context) { globalController.SetGlobalState(c, app) })
	protectedRouter.GET("/state", func(c *gin.Context) { globalController.GetState(c, app) })
	protectedRouter.POST("/global/tuition_type", func(c *gin.Context) { globalController.SetTuitionType(c, app) })
	protectedRouter.GET("/global/tuition_type", func(c *gin.Context) { globalController.GetTuitionType(c, app) })
	// auth routes
	protectedRouter.GET("/logout", func(c *gin.Context) { authController.Logout(c, app) })
	protectedRouter.POST("/refresh-token", func(c *gin.Context) { authController.RefreshToken(c, app) })

	// user routes
	protectedRouter.GET("/user/info", func(c *gin.Context) { usersController.GetUserInfo(c, app) })
	protectedRouter.GET("/user/list", func(c *gin.Context) { usersController.ListUsers(c, app) })
	protectedRouter.PATCH("/user/change-password", func(c *gin.Context) { usersController.ChangePassUser(c, app) })
	protectedRouter.PUT("/user/update", func(c *gin.Context) { usersController.UpdateUser(c, app) })

	// user - course routes
	protectedRouter.GET("/user/course/list", func(c *gin.Context) { usersCourseController.UserListCourse(c, app) })
	protectedRouter.POST("/user/course/register", func(c *gin.Context) { usersCourseController.UserRegisterCourse(c, app) })
	protectedRouter.POST("/user/course/unregister", func(c *gin.Context) { usersCourseController.UserUnRegisterCourse(c, app) })
	protectedRouter.POST("/user/course/unregister_admin", func(c *gin.Context) { usersCourseController.UserUnRegisterAdminCourse(c, app) })
	protectedRouter.DELETE("/user/course/delete", func(c *gin.Context) { usersCourseController.UserDeleteCourse(c, app) })

	// course routes
	protectedRouter.POST("/course/create", func(c *gin.Context) { coursesController.CreateCourse(c, app) })
	protectedRouter.DELETE("/course/delete/:course_id", func(c *gin.Context) { coursesController.DeleteCourse(c, app) })
	protectedRouter.PUT("/course/edit", func(c *gin.Context) { coursesController.EditCourse(c, app) })
	protectedRouter.PUT("/course/confirm", func(c *gin.Context) { coursesController.ConfirmCourse(c,app)})

	// department routes
	protectedRouter.POST("/department/create", func(c *gin.Context) { departmentsController.CreateDepartment(c, app) })
	protectedRouter.PUT("/department/update", func(c *gin.Context) { departmentsController.UpdateDepartment(c, app) })

	//lecturer routes
	protectedRouter.GET("/lecturer/course/list", func(c *gin.Context) { lecturerCourseController.ListLecturerCourses(c, app) })
	protectedRouter.GET("/lecturer/course/enroller/list", func(c *gin.Context) { lecturerCourseController.ListCourseEnroller(c, app) })
	protectedRouter.GET("/lecturer/course/register/list", func(c *gin.Context) { lecturerCourseController.ListLecturerRegisterCourses(c, app) })

	protectedRouter.POST("/user/credit", func(c *gin.Context) { usersCourseController.GetCredit(c, app) })
	protectedRouter.POST("/tuition/cal_tuition", func(c *gin.Context) { tuitionController.CalTuition(c, app) })
	protectedRouter.POST("/tuition/create_tuition", func(c *gin.Context) { tuitionController.CreateTuition(c, app)})
	protectedRouter.GET("/tuition/list", func(c *gin.Context) { tuitionController.ListTuition(c, app) })
	protectedRouter.POST("/tuition/pay", func(c *gin.Context) { tuitionController.PayTuition(c, app) })

}
