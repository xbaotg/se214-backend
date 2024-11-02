package auth

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

// Logout user
// @Summary Logout user
// @Description Logout user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /logout [post]
func Logout(c *gin.Context, app *bootstrap.App) {
	sess, exists := c.Get("session")

	if !exists {
		app.Logger.Error().Msg("Session not found")
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// get user from refresh token
	session := sess.(models.Session)

	// revoke refresh token
	if err := app.DB.Model(&session).Updates(models.Session{IsActive: false}).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Logout success", nil)
}
