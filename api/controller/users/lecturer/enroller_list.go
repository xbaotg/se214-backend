package lecturer;

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// List course enroller
// @Summary List course enroller
// @Description List course enroller
// @Tags Lecturer
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param course_id query string false "Course ID"
// @Router /lecturer/course/enroller/list [get]
func ListCourseEnroller(c *gin.Context, app *bootstrap.App) {
	courseID := c.Query("course_id")
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// if !slices.Contains([]models.Role{models.RoleAdmin, models.RoleUser}, user.UserRole) {
	// 	internal.Respond(c, 403, false, "Forbidden", nil)
	// 	return
	// }
	if user.UserRole == models.RoleUser {
		internal.Respond(c, 403, false, "Forbidden", nil)
		return
	}

	course := models.Course{}
	if err := app.DB.Where("id = ?", courseID).First(&course).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	users := []models.User{}
	if err := app.DB.Select("users.*").Table("courses").Joins(
		"JOIN registered_courses ON courses.id = registered_courses.course_id",
		).Joins(
			"JOIN users ON registered_courses.user_id = users.id",
		).Where("courses.id = ?", courseID).Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			internal.Respond(c, 404, false, "Course enroller not found", nil)
			return
		}
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "List course enroller successfully", users)
}