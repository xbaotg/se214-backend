package global

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

type Stats struct {
	TotalUsers int64 `json:"total_users"`
	TotalTeachers int64 `json:"total_teachers"`
	TotalStudents int64 `json:"total_students"`
	TotalCourses int64 `json:"total_courses"`
	TotalSubjects int64 `json:"total_subjects"`
	TotalStudentsRegistered int64 `json:"total_students_registered"`
	TotalStudentPaid int64 `json:"total_student_paid"`
	TotalCourseRequest int64 `json:"total_course_request"`
	TotalMoney int64 `json:"total_money"`
}


// GetStats godoc
// @Summary Get stats
// @Description Get stats
// @Tags Global
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /global/stats [get]
func GetStats(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	// get user info
	user := models.User{
		ID: session.UserID,
	}
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	var stats Stats

	if err := app.DB.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}

	if err := app.DB.Model(&models.User{}).Where("user_role = ?", models.RoleLecturer).Count(&stats.TotalTeachers).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}

	if err := app.DB.Model(&models.User{}).Where("user_role = ?", models.RoleUser).Count(&stats.TotalStudents).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}

	if err := app.DB.Model(&models.Course{}).Count(&stats.TotalCourses).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}

	if err := app.DB.Model(&models.AllCourses{}).Count(&stats.TotalSubjects).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}

	if err := app.DB.Table("registered_courses").Select("DISTINCT user_id").Count(&stats.TotalStudentsRegistered).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}


	if err := app.DB.Table("courses").Where("confirmed = ?", false).Count(&stats.TotalCourseRequest).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}

	if err := app.DB.Table("tuitions").Select("SUM(cost)").Where("status = ? and year = ? and semester = ?", 
		models.TuStatusPaid, app.Config.CurrentYear, app.Config.CurrentSemester).Count(&stats.TotalMoney).Error; err != nil {
	}

	if err := app.DB.Table("tuitions").Select("DISTINCT user_id").Where("status = ? and year = ? and semester = ?", 
		models.TuStatusPaid, app.Config.CurrentYear, app.Config.CurrentSemester).Count(&stats.TotalStudentPaid).Error; err != nil {
	}

	internal.Respond(c, 200, true, "Thống kê", stats)
}
