package tuition

import (
	"be/internal"
	"be/bootstrap"
	"be/models"
	"be/api/schemas"
	// "net/http"
	// "slices"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CalTuitionAdminRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Year int `json:"year" binding:"required"`
	Semester int `json:"semester" binding:"required"`
}

// GetTuition godoc
// @Summary Calculate tuition
// @Description Calculate tuition
// @Tags Tuition
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param CalTuitionAdminRequest body CalTuitionAdminRequest true "CalTuitionAdminRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /tuition/get_tuition [post]
func GetTuition(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	// if app.State == bootstrap.ACTIVE {
	// 	internal.Respond(c, 403, false, "Cannot calculate tuition in this state", nil)
	// 	return
	// }
	var req CalTuitionAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, "Hãy điền đầy đủ thông tin", nil)
		return
	}

	// get user info
	user := models.User{
		ID: session.UserID,
	}
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if user.UserRole == models.RoleAdmin {
		user.ID = req.UserID
	} else if user.UserRole == models.RoleUser {
		user.ID = session.UserID
	} else {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	var courses []models.Course
	var credit int32 = 0
	if err := app.DB.Select("courses.*").Table(
		"registered_courses").Joins("JOIN courses ON courses.id = registered_courses.course_id").Where(
		"registered_courses.user_id = ? AND registered_courses.course_year = ? AND registered_courses.course_semester = ?",
		user.ID, req.Year, req.Semester).Find(&courses).Error; err != nil {
		

		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}	

	coursesResponse := []schemas.CreateCourseResponse{}
	if len(courses) == 0 {
		internal.Respond(c, 200, true, "Người dùng chưa đăng ký môn học nào", gin.H{"tuition": -1})
		return
	} else {
		for _, course := range courses {
			credit += course.CourseCredit
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
	}

	tuition := models.Tuition{}
	if err := app.DB.Where("user_id = ? AND year = ? AND semester = ?", user.ID, req.Year, req.Semester).First(&tuition).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}
	internal.Respond(c, 200, true, "Tuition", gin.H{
		"courses": coursesResponse,
		"credit": credit,
		"tuition": tuition.Tuition,
	})
}
