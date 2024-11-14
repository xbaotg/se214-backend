package coursesmanagament

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserListCourseRequest struct {
	CourseYear     int32 `json:"course_year"`
	CourseSemester int32 `json:"course_semester"`
}

func UserListCourse(c *gin.Context, app *bootstrap.App) {
	req := UserListCourseRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
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

	// find courses
	registeredCourses := []models.RegisteredCourse{}
	if err := app.DB.Where("user_id = ? AND course_year = ? AND course_semester = ?", user.ID, req.CourseYear, req.CourseSemester).Find(&registeredCourses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Respond(c, 404, false, "No courses found", nil)
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
	}

	// map registered courses to courses
	courses := []models.Course{}
	for _, registeredCourse := range registeredCourses {
		course := models.Course{}
		if err := app.DB.First(&course, "id = ?", registeredCourse.CourseID).Error; err != nil {
			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}

		courses = append(courses, course)
	}

	internal.Respond(c, 200, true, "Courses found", courses)
}
