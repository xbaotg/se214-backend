package auth

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	SessionID             uuid.UUID   `json:"session_id"`
	UserRole              models.Role `json:"user_role"`
	AccessToken           string      `json:"access_token"`
	AccessTokenExpiresIn  time.Time   `json:"access_token_expires_in"`
	RefreshToken          string      `json:"refresh_token"`
	RefreshTokenExpiresIn time.Time   `json:"refresh_token_expires_in"`
}

// Login user
// @Summary Login user
// @Description Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param LoginRequest body LoginRequest true "LoginRequest"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /login [post]
func Login(c *gin.Context, app *bootstrap.App) {
	req := LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	// user, err := app.DB.GetUserByUsername(c, req.Username)
	user := models.User{}
	if err := app.DB.First(&user, "username = ?", req.Username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Respond(c, 404, false, "User not found", nil)
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
	}

	if internal.CheckPassword(req.Password, user.Password) != nil {
		internal.Respond(c, 400, false, "Invalid password", nil)
		return
	}

	// create access token
	accessToken, accessPayload, err := app.TokenMaker.CreateToken(user.Username, time.Second*time.Duration(app.Config.JWTExpire))

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// create refresh token
	refreshToken, refreshPayload, err := app.TokenMaker.CreateToken(
		user.Username,
		time.Second*time.Duration(app.Config.JWTRefreshExpire),
	)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	session := models.Session{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresIn:    refreshPayload.ExpiredAt,
		IsActive:     true,
		CreatedAt:    internal.GetCurrentTime(),
		UpdatedAt:    internal.GetCurrentTime(),
	}

	if err := app.DB.Create(&session).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, http.StatusOK, true, "Login success", LoginResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		UserRole:              user.UserRole,
		AccessTokenExpiresIn:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: refreshPayload.ExpiredAt,
	})
}
