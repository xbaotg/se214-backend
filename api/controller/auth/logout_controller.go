package controller

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"
	"be/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout user
// @Summary Logout user
// @Description Logout user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /logout [post]
func Logout(c *gin.Context, app *bootstrap.App) {
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

	// revoke refresh token
	err := app.DB.RevolveSession(c, sqlc.RevolveSessionParams{
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
