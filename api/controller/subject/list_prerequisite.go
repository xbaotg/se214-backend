package subject

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

type PrerequisiteResponse struct {
	PrerequisiteID string `json:"prerequisite_id"`
	CourseFullname string `json:"course_fullname"`
}
// ListPrerequisite docs
// @Summary List prerequisite
// @Description List prerequisite
// @Tags Subject
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param course_name query string true "Course name"
// @Router /subject/prerequisite [get]
func ListPrerequisite(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	// validate user
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 404, false, "User not found", nil)
		return
	}

	courseName := c.Query("course_name")

	// if user.UserRole != models.RoleAdmin {
	// 	internal.Respond(c, 403, false, "Permission denied", nil)
	// 	return
	// }

	results := []PrerequisiteResponse{}
	if err:= app.DB.Table("prerequisite_courses").Joins(
		"LEFT JOIN all_courses ON prerequisite_courses.prerequisite_id = all_courses.course_name").Where(
			"prerequisite_courses.course_id = ?", courseName).Select(
				"prerequisite_courses.prerequisite_id, all_courses.course_fullname").Scan(&results).Error; err != nil {

		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Subjects fetched successfully", results)
}

