package controller

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserInfoResponse struct {
	UserID       uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	UserEmail    string    `json:"email"`
	UserFullname string    `json:"user_fullname"`
	UserRole     sqlc.Role `json:"user_role"`
	Year         int32     `json:"year"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Get user info
// @Summary Get user info
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /user-info [post]
func GetUserInfo(c *gin.Context, app *bootstrap.App) {
	sess, exists := c.Get("session")
	if !exists {
		app.Logger.Error().Msg("Session not found")
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	session := sess.(sqlc.Session)

	user, err := app.DB.GetUserById(c, session.UserID)
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
		Message: "Get user info successfully",
		Data: UserInfoResponse{
			UserID:       user.UserID,
			Username:     user.Username,
			UserEmail:    user.UserEmail,
			UserFullname: user.UserFullname,
			UserRole:     user.UserRole,
			Year:         user.Year,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		},
	})
}
