package subject

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	// "fmt"

	"github.com/gin-gonic/gin"
)

// DeleteSubject docs
// @Summary Delete subject
// @Description Delete subject
// @Tags Subject
// @Accept json
// @Produce json
// @Param course_name query string true "Course name"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /subject/delete [delete]
func DeleteSubject(c *gin.Context, app *bootstrap.App) {
	// if app.State != bootstrap.SETUP {
	// 	internal.Respond(c, 403, false, fmt.Sprintf("Máy chủ không ở trạng thái SETUP, trạng thái hiện tại là %s", app.State), nil)
	// 	return
	// }

	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	// validate user
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 404, false, "Người dùng không tồn tại", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}
	courseName := c.Query("course_name")
	// check if course already exists
	course := models.AllCourses{CourseName: courseName}
	if err := app.DB.Table("all_courses").Where(course).Delete(&course).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Xóa môn học thành công", course)
}

