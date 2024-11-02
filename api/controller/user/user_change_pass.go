package user

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"
	"be/model"

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
		c.JSON(400, model.Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	if req.OldPassword == req.NewPassword {
		c.JSON(400, model.Response{
			Status:  false,
			Message: "New password must be different from old password",
		})
		return
	}

	user, err := app.DB.GetUserById(c, session.UserID)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	if internal.CheckPassword(req.OldPassword, user.Password) != nil {
		c.JSON(400, model.Response{
			Status:  false,
			Message: "Old password is incorrect",
		})
		return
	}

	hashedPassword, err := internal.HashPassword(req.NewPassword)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	_, err = app.DB.UpdatPassword(c, sqlc.UpdatPasswordParams{
		UserID:    session.UserID,
		Password:  hashedPassword,
		UpdatedAt: internal.GetCurrentTime(),
	})

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(200, model.Response{
		Status:  true,
		Message: "Change password success",
	})

}
