package controller

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"
	"be/model"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken           string
	AccessTokenExpiresIn  time.Time
	RefreshToken          string
	RefreshTokenExpiresIn time.Time
}

func RefreshToken(c *gin.Context, app *bootstrap.App) {
	req := RefreshTokenRequest{}

	// validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	// validate refresh token
	refreshPayload, err := app.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, model.Response{
			Status:  false,
			Message: "Invalid token",
		})
		return
	}

	// get user from refresh token
	session, err := app.DB.GetSessionBySessionId(c, refreshPayload.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, model.Response{
				Status:  false,
				Message: "Session not found",
			})
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(500, model.Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	// check if refresh token is valid
	if time.Now().After(session.ExpiresIn) {
		c.JSON(401, model.Response{
			Status:  false,
			Message: "Refresh token expired",
		})
		return
	}

	if (!session.IsActive) || (session.RefreshToken != req.RefreshToken) {
		c.JSON(401, model.Response{
			Status:  false,
			Message: "Invalid token",
		})
		return
	}

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
