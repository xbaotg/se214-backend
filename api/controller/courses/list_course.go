package courses

import (
	"be/bootstrap"

	"github.com/gin-gonic/gin"
)

type ListCourseRequest struct {
}

func ListCourses(c *gin.Context, app *bootstrap.App) {
	// sess, _ := c.Get("session")
	// session := sess.(models.Session)

	// // request validation
	// req := ListCourseRequest{}
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	internal.Respond(c, 400, false, err.Error(), nil)
	// 	return
	// }

	// courses, err := app.ListCourses(session)
	// if err != nil {
	// 	internal.Respond(c, 500, false, err.Error(), nil)
	// 	return
	// }

	// internal.Respond(c, 200, true, "List courses success", courses)
}
