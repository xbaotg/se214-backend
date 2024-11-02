package auth

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	SessionID             uuid.UUID
	AccessToken           string
	AccessTokenExpiresIn  time.Time
	RefreshToken          string
	RefreshTokenExpiresIn time.Time
}

// Login user
// @Summary Login user
// @Description Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param LoginRequest body controller.LoginRequest true "LoginRequest"
// @Success 200 {object} controller.LoginResponse
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /login [post]
func Login(c *gin.Context, app *bootstrap.App) {
	req := LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	user, err := app.DB.GetUserByUsername(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
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

	session, err := app.DB.CreateSession(c, sqlc.CreateSessionParams{
		SessionID:    refreshPayload.ID,
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		ExpiresIn:    refreshPayload.ExpiredAt,
		IsActive:     true,
		CreatedAt:    internal.GetCurrentTime(),
		UpdatedAt:    internal.GetCurrentTime(),
	})

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, http.StatusOK, true, "Login success", LoginResponse{
		SessionID:             session.SessionID,
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: refreshPayload.ExpiredAt,
	})
}
