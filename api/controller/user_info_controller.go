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

type UserInfoRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

// Get user info
// @Summary Get user info
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Param UserInfoRequest body controller.UserInfoRequest true "UserInfoRequest"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /user-info [post]
func GetUserInfo(c *gin.Context, app *bootstrap.App) {
	var req UserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	payload, err := app.TokenMaker.VerifyToken(req.AccessToken)
	if err != nil {
		c.JSON(401, model.Response{
			Status:  false,
			Message: "Invalid token",
		})
		return
	}

	user, err := app.DB.GetUserByUsername(c, payload.Username)
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
