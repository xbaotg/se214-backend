package courses

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteCourseRequest struct {
	CourseID uuid.UUID `json:"course_id" binding:"required"`
}

// @Summary Delete course
// @Description Delete course
// @Tags Courses
// @Accept json
// @Produce json
// @Param course_id body string true "Course ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /course/delete [delete]
func DeleteCourse(c *gin.Context, app *bootstrap.App) {
	// validate request
	req := DeleteCourseRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	// check if user exists
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, err.Error(), nil)
		return
	}

	// check if user is admin
	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Permission denied", nil)
		return
	}

	course := models.Course{ID: req.CourseID}

	// check if course exists
	if err := app.DB.First(&course).Error; err != nil {
		internal.Respond(c, 404, false, "Course not found", nil)
		return
	}

	// delete course
	if err := app.DB.Delete(&course).Error; err != nil {
		internal.Respond(c, 500, false, err.Error(), nil)
		return
	}

	internal.Respond(c, 200, true, "Delete course success", nil)
}
