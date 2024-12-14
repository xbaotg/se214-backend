package coursesmanagament

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type CourseDeleteRequest struct {
// 	CourseID uuid.UUID `json:"course_id" binding:"required"`
// }

// CourseDelete docs
// @Summary Delete course
// @Description Delete course
// @Tags User
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /user/course/delete [delete]
func UserDeleteCourse(c *gin.Context, app *bootstrap.App) {
	// validate request
	// req := UserRegisterCourseRequest{}
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	internal.Respond(c, 400, false, err.Error(), nil)
	// 	return
	// }
	if app.State != bootstrap.ACTIVE {
		internal.Respond(c, 403, false, fmt.Sprintf("Server is not in active state, current state is %s", app.State), nil)
		return
	}

	courseID_ := c.Param("course_id")
	courseID, err := uuid.Parse(courseID_)
	if err != nil {
		internal.Respond(c, 400, false, "Invalid course ID", nil)
		return
	}

	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	// validate user
	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// validate course
	course := models.Course{ID: courseID}
	if err := app.DB.First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Respond(c, 404, false, "Course not found", nil)
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
	}

	// check if user registered the course
	registeredCourse := models.RegisteredCourse{UserID: user.ID, CourseID: course.ID}
	if err := app.DB.First(&registeredCourse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Respond(c, 404, false, "User has not registered the course", nil)
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
	}

	// delete the course
	if err := app.DB.Delete(&registeredCourse).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Course deleted successfully", nil)
}
