package controller

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"
	"be/model"
	"time"

	"github.com/gin-gonic/gin"
)

type RefreshTokenResponse struct {
	AccessToken           string
	AccessTokenExpiresIn  time.Time
	RefreshToken          string
	RefreshTokenExpiresIn time.Time
}

// Refresh token
// @Summary Refresh token
// @Description Refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} controller.RefreshTokenResponse
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /refresh-token [post]
func RefreshToken(c *gin.Context, app *bootstrap.App) {
	sess, exists := c.Get("session")

	if !exists {
		app.Logger.Error().Msg("Session not found")

		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	// get user from refresh token
	session := sess.(sqlc.Session)

	// generate new access token
	accessToken, AccessTokenPayload, err := app.TokenMaker.CreateToken(session.UserID.String(), time.Second*time.Duration(app.Config.JWTExpire))
	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	// generate new refresh token
	refreshToken, refreshPayload, err := app.TokenMaker.CreateToken(session.UserID.String(), time.Second*time.Duration(app.Config.JWTRefreshExpire))
	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	// update session
	_, err = app.DB.CreateSession(c, sqlc.CreateSessionParams{
		SessionID:    refreshPayload.ID,
		UserID:       session.UserID,
		RefreshToken: refreshToken,
		ExpiresIn:    refreshPayload.ExpiredAt,
		IsActive:     true,
		CreatedAt:    internal.GetCurrentTime(),
		UpdatedAt:    internal.GetCurrentTime(),
	})
	if err != nil {
		app.Logger.Error().Err(err).Msg("Failed to create session: " + err.Error())
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	// revoke old refresh token
	err = app.DB.RevolveSession(c, sqlc.RevolveSessionParams{
		SessionID: session.SessionID,
		UpdatedAt: internal.GetCurrentTime(),
		IsActive:  false,
	})

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(200, RefreshTokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  AccessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: refreshPayload.ExpiredAt,
	})
}
