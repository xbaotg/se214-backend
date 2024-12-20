package courses

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"be/api/schemas"

	"github.com/gin-gonic/gin"
)

type ListCourseRequest struct {
}

// List Courses docs
// @Summary List Courses
// @Description List all courses
// @Tags Courses
// @Accept json
// @Produce json
// @Success 200 {object} schemas.CreateCourseResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /course/list [get]
func ListCourses(c *gin.Context, app *bootstrap.App) {
	courses := []models.Course{}
	if err := app.DB.Where("confirmed = ?", true).Find(&courses).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	coursesResponse := []schemas.CreateCourseResponse{}
	for _, course := range courses {
		coursesResponse = append(coursesResponse, schemas.CreateCourseResponse{
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

	internal.Respond(c, 200, true, "Lấy danh sách khóa học thành công", coursesResponse)
}
