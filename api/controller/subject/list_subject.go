package subject

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

type SubjectsResponse struct {
	CourseName string `json:"course_name"`
	CourseFullname	string `json:"course_fullname"`
	DepartmentCode string `json:"department_code"`
}

// ListSubject docs
// @Summary List subject
// @Description List subject
// @Tags Subject
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /subject/list [get]
func ListSubject(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	// validate user
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 404, false, "Người dùng không tồn tại", nil)
		return
	}

	// if user.UserRole != models.RoleAdmin {
	// 	internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
	// 	return
	// }

	results := []SubjectsResponse{}
	if err:= app.DB.Table("all_courses").Joins("LEFT JOIN courses ON all_courses.course_name = courses.course_name").Joins(
		"LEFT JOIN departments ON courses.department_id = departments.id").Where("all_courses.status = ?", true).Select(
			"all_courses.course_name, all_courses.course_fullname, departments.department_code").Scan(&results).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}
	internal.Respond(c, 200, true, "Danh sách môn học", results)
}

