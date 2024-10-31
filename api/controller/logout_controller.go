package controller

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"
	"be/model"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Logout user
// @Summary Logout user
// @Description Logout user
// @Tags Auth
// @Accept json
// @Produce json
// @Param LogoutRequest body controller.LogoutRequest true "LogoutRequest"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /logout [post]
func Logout(c *gin.Context, app *bootstrap.App) {
	req := LogoutRequest{}

	// validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
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

	// get session from refresh token
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

	// revoke refresh token
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

	c.JSON(http.StatusOK, model.Response{
		Status:  true,
		Message: "Logout success",
	})
}
