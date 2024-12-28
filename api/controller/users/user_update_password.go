package users

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

type ChangePassRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Change password
// @Summary Change password
// @Description Change password
// @Tags User
// @Accept json
// @Produce json
// @Param ChangePassRequest body ChangePassRequest true "ChangePassRequest"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/change-pass [post]
func ChangePassUser(c *gin.Context, app *bootstrap.App) {
	// validate request
	req := ChangePassRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	if req.OldPassword == req.NewPassword {
		internal.Respond(c, 400, false, "Mật khẩu mới không được trùng với mật khẩu cũ", nil)
		return
	}

	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if internal.CheckPassword(req.OldPassword, user.Password) != nil {
		internal.Respond(c, 400, false, "Mật khẩu cũ không đúng", nil)
		return
	}

	hashedPassword, err := internal.HashPassword(req.NewPassword)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if err := app.DB.Model(&user).Updates(models.User{Password: hashedPassword}).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Đổi mật khẩu thành công", nil)
}
