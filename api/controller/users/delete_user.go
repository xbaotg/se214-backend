package users

import (
	"be/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"be/models"
	"be/internal"	
)


// Delete user
// @Summary Delete user
// @Description Delete user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id query string true "User ID"
// @Router /user/delete [delete]
func DeleteUser(c *gin.Context, app *bootstrap.App) {
	app.Logger.Info().Msg("Update user")

	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	id := c.Query("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		internal.Respond(c, 400, false, "ID không hợp lệ", nil)
		return
	}

	if err := app.DB.Table("users").Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Xóa người dùng thành công", nil)
}