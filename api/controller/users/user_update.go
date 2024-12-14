package users

import (
	"be/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"be/models"
	"be/internal"	
)


type UpdateUserRequest struct {
	UserID 	 uuid.UUID `json:"id" binding:"required"`
	UserFullname string      `json:"user_fullname"`
	UserRole     models.Role `json:"user_role"`
	Year         int32       `json:"year"`

}

// Update user
// @Summary Update user
// @Description Update user
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUserRequest body UpdateUserRequest true "UpdateUserRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /user/update [post]
func UpdateUser(c *gin.Context, app *bootstrap.App) {
	req := UpdateUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Forbidden", nil)
		return
	}

	updateUser := models.User{ID: req.UserID}

	if err := app.DB.First(&updateUser).Error; err != nil {
		internal.Respond(c, 404, false, "User not found", nil)
		return
	}

	updateUser.UserFullname = req.UserFullname
	updateUser.UserRole = req.UserRole
	updateUser.Year = req.Year

	if err := app.DB.Save(&updateUser).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Update user successfully", nil)
}