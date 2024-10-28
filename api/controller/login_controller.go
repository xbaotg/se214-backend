package controller

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"
	. "be/model"
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

func Login(c *gin.Context, app *bootstrap.App) {
	req := LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	user, err := app.DB.GetUserByUsername(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, Response{
				Status:  false,
				Message: "User not found",
			})
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(500, Response{
			Status:  false,
			Message: "Internal server error",
		})
	}

	if internal.CheckPassword(req.Password, user.Password) != nil {
		c.JSON(400, Response{
			Status:  false,
			Message: "Invalid password",
		})
		return
	}

	// create access token
	accessToken, accessPayload, err := app.TokenMaker.CreateToken(user.Username, time.Second*time.Duration(app.Config.JWTExpire))

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Status:  false,
			Message: "Internal server error",
		})
		return
	}

	// create refresh token
	refreshToken, refreshPayload, err := app.TokenMaker.CreateToken(
		user.Username,
		time.Second*time.Duration(app.Config.JWTRefreshExpire),
	)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())

		c.JSON(http.StatusInternalServerError, Response{
			Status:  false,
			Message: err.Error(),
		})

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

		c.JSON(http.StatusInternalServerError, Response{
			Status:  false,
			Message: "Internal server error",
		})

		return
	}

	c.JSON(200, Response{
		Status:  true,
		Message: "Login success",
		Data: LoginResponse{
			SessionID:             session.SessionID,
			AccessToken:           accessToken,
			AccessTokenExpiresIn:  accessPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresIn: refreshPayload.ExpiredAt,
		},
	})
}
