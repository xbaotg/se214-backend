package user

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"

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
// @Param ChangePassRequest body user.ChangePassRequest true "ChangePassRequest"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /user/change-pass [post]
func ChangePassUser(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(sqlc.Session)

	// validate request
	req := ChangePassRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	if req.OldPassword == req.NewPassword {
		internal.Respond(c, 400, false, "New password must be different from old password", nil)
		return
	}

	user, err := app.DB.GetUserById(c, session.UserID)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	if internal.CheckPassword(req.OldPassword, user.Password) != nil {
		internal.Respond(c, 400, false, "Old password is incorrect", nil)
		return
	}

	hashedPassword, err := internal.HashPassword(req.NewPassword)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	_, err = app.DB.UpdatPassword(c, sqlc.UpdatPasswordParams{
		UserID:    session.UserID,
		Password:  hashedPassword,
		UpdatedAt: internal.GetCurrentTime(),
	})

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Change password success", nil)
}
