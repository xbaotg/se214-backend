package courses

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

type ListCourseRequest struct {
}

func ListCourses(c *gin.Context, app *bootstrap.App) {
	courses := []models.Course{}
	if err := app.DB.Find(&courses).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Courses found", courses)
}
