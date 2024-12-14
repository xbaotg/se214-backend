package coursesmanagament

import (
	"be/api/schemas"
	"be/bootstrap"
	"be/internal"
	"be/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// type UserListCourseRequest struct {
// 	CourseYear     int32 `json:"course_year"`
// 	CourseSemester int32 `json:"course_semester"`
// }

// UserListCourse docs
// @Summary User list course
// @Description User list course
// @Tags User
// @Accept json
// @Produce json
// @Param course_year query int true "Course year"
// @Param course_semester query int true "Course semester"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /user/course/list [GET]
func UserListCourse(c *gin.Context, app *bootstrap.App) {
	// req := UserListCourseRequest{}
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	internal.Respond(c, 400, false, err.Error(), nil)
	// 	return
	// }
	courseYear := c.Query("course_year")
	courseSemester := c.Query("course_semester")

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
	if err := app.DB.Where("user_id = ? AND course_year = ? AND course_semester = ?", user.ID, courseYear, courseSemester).Find(&registeredCourses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Respond(c, 404, false, "No courses found", nil)
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
	}

	// map registered courses to courses
	courses := []schemas.CreateCourseResponse{}
	for _, registeredCourse := range registeredCourses {
		course := models.Course{}
		if err := app.DB.First(&course, "id = ? and confirmed=?", registeredCourse.CourseID, true).Error; err != nil {
			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}

		teacher := models.User{}
		if err := app.DB.First(&teacher, "id = ?", course.CourseTeacherID).Error; err != nil {
			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}

		courses = append(courses, schemas.CreateCourseResponse{
			ID:                course.ID,
			CourseTeacherID:   course.CourseTeacherID,
			CourseTeacherName: teacher.UserFullname,
			CourseDepartment:  course.DepartmentID,
			CourseName:        course.CourseName,
			CourseFullname:    course.CourseFullname,
			CourseCredit:      course.CourseCredit,
			CourseYear:        course.CourseYear,
			CourseSemester:    course.CourseSemester,
			CourseStartShift:  course.CourseStartShift,
			CourseEndShift:    course.CourseEndShift,
			CourseDay:         course.CourseDay,
			MaxEnroll:         course.MaxEnroller,
			CurrentEnroll:     course.CurrentEnroller,
			CourseRoom:        course.CourseRoom,
		})
	}

	internal.Respond(c, 200, true, "Courses found", courses)
}
