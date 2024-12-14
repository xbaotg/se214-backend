package lecturer

import (
	"be/api/schemas"
	"be/bootstrap"
	"be/internal"
	"be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// List lecturer courses
// @Summary List lecturer courses
// @Description List lecturer courses
// @Tags Lecturer
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /lecturer/course/list [get]
func ListLecturerCourses(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// if user.Role != models.Lecturer {
	// 	internal.Respond(c, 403, false, "Forbidden", nil)
	// 	return
	// }

	courses := []models.Course{}
	if user.UserRole == models.RoleLecturer {
		if err := app.DB.Where("course_teacher_id = ? and confirmed = ?", user.ID, true).Find(&courses).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				internal.Respond(c, 404, false, "Course not found", nil)
				return
			}
			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}
	} else if user.UserRole == models.RoleAdmin {
		if err := app.DB.Where("confirmed = ?", true).Find(&courses).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				internal.Respond(c, 404, false, "Course not found", nil)
				return
			}
			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}
	} else {
		internal.Respond(c, 403, false, "Forbidden", nil)
		return
	}

	var courseResponse []schemas.CreateCourseResponse
	for _, course := range courses {
		courseResponse = append(courseResponse, schemas.CreateCourseResponse{
			ID:               course.ID,
			CourseTeacherID:  course.CourseTeacherID,
			CourseDepartment: course.DepartmentID,
			CourseName:       course.CourseName,
			CourseFullname:   course.CourseFullname,
			CourseCredit:     course.CourseCredit,
			CourseYear:       course.CourseYear,
			CourseSemester:   course.CourseSemester,
			CourseStartShift: course.CourseStartShift,
			CourseEndShift:   course.CourseEndShift,
			CourseDay:        course.CourseDay,
			MaxEnroll:        course.MaxEnroller,
			CurrentEnroll:    course.CurrentEnroller,
			CourseRoom:       course.CourseRoom,
		})
	}

	internal.Respond(c, 200, true, "List lecturer courses successfully", courseResponse)
}