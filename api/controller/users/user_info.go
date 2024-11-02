package users

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserInfoResponse struct {
	UserID       uuid.UUID   `json:"id"`
	Username     string      `json:"username"`
	UserEmail    string      `json:"email"`
	UserFullname string      `json:"user_fullname"`
	UserRole     models.Role `json:"user_role"`
	Year         int32       `json:"year"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// Get user info
// @Summary Get user info
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/info [post]
func GetUserInfo(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Get user info successfully", UserInfoResponse{
		UserID:       user.ID,
		Username:     user.Username,
		UserEmail:    user.UserEmail,
		UserFullname: user.UserFullname,
		UserRole:     user.UserRole,
		Year:         user.Year,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})
}
