package tuition

import (
	"be/internal"
	"be/bootstrap"
	"be/models"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Delete tuition godoc
// @Summary Delete tuition
// @Description Delete tuition
// @Tags Tuition
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param id query string true "Tuition ID"
// @Router /tuition/delete [delete]
func DeleteTuition(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	// if app.State == bootstrap.ACTIVE {
	// 	internal.Respond(c, 403, false, "Cannot calculate tuition in this state", nil)
	// 	return
	// }
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

	id := c.Query("id")
	tuitionID, err := uuid.Parse(id)
	if err != nil {
		internal.Respond(c, 400, false, "ID không hợp lệ", nil)
		return
	}

	if err := app.DB.Table("tuitions").Where("id = ?", tuitionID).Delete(&models.Tuition{}).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Xóa học phí thành công", nil)
}
