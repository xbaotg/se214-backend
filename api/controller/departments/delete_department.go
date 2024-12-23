package departments

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

// @Summary Delete department
// @Description Delete department
// @Tags Department
// @Accept json
// @Produce json
// @Success 200 {object} models.Department
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param department_code query string true "Department code"
// @Router /department/delete [delete]
func DeleteDepartment(c *gin.Context, app *bootstrap.App) {

	departmentCode := c.Query("department_code")

	// get current user from session
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	// check if user is admin
	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	var department models.Department
	if err := app.DB.Table("departments").Where("department_code = ?", departmentCode).First(&department).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	var courses []models.Course
	if err := app.DB.Table("courses").Where("department_id = ?", department.ID).Find(&courses).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if len(courses) > 0 {
		internal.Respond(c, 400, false, "Không thể xóa khoa đã có môn học", nil)
		return
	}

	if err := app.DB.Delete(&department).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Xóa khoa thành công", nil)
}
