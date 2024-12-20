package subject

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	// "fmt"

	"github.com/gin-gonic/gin"
)

type SubjectCreateRequest struct {
	CourseName string `json:"course_name" binding:"required"`
	CourseFullname	string `json:"course_fullname" binding:"required"`
}

// CreatteSubject docs
// @Summary Create subject
// @Description Create subject
// @Tags Subject
// @Accept json
// @Produce json
// @Param SubjectCreateRequest body SubjectCreateRequest true "SubjectCreateRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /subject/create [post]
func CreateSubject(c *gin.Context, app *bootstrap.App) {
	// if app.State != bootstrap.SETUP {
	// 	internal.Respond(c, 403, false, fmt.Sprintf("Server is not in SETUP state, current state is %s", app.State), nil)
	// 	return
	// }
	// validate request
	req := SubjectCreateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	// validate user
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 404, false, "User not found", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Permission denied", nil)
		return
	}

	// check if course already exists
	course := models.AllCourses{CourseName: req.CourseName, CourseFullname: req.CourseFullname}
	if res := app.DB.Table("all_courses").Where("course_name = ?", req.CourseName).FirstOrCreate(&course); res.Error != nil {
		if res.RowsAffected == 0 {
			app.Logger.Info().Msg("Course already exists")
			if err := app.DB.Table("all_courses").Where(course).Update("status", true).Error; err != nil {
				app.Logger.Error().Err(err).Msg(err.Error())
				internal.Respond(c, 500, false, "Internal server error", nil)
				return
			}

			if req.CourseFullname != "" {
				if err := app.DB.Table("all_courses").Where(course).Update("course_fullname", req.CourseFullname).Error; err != nil {
					app.Logger.Error().Err(err).Msg(err.Error())
					internal.Respond(c, 500, false, "Internal server error", nil)
					return
				}
			}
			internal.Respond(c, 200, true, "Course updated successfully", course)
		}
		app.Logger.Error().Err(res.Error).Msg(res.Error.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Course created successfully", course)
}

