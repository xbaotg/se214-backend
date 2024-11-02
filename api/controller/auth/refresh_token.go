package auth

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
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
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /refresh-token [post]
func RefreshToken(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	// generate new access token
	accessToken, AccessTokenPayload, err := app.TokenMaker.CreateToken(session.UserID.String(), time.Second*time.Duration(app.Config.JWTExpire))
	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// generate new refresh token
	refreshToken, refreshPayload, err := app.TokenMaker.CreateToken(session.UserID.String(), time.Second*time.Duration(app.Config.JWTRefreshExpire))
	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// update session
	updatedSession := models.Session{
		ID:           refreshPayload.ID,
		UserID:       session.UserID,
		RefreshToken: refreshToken,
		ExpiresIn:    refreshPayload.ExpiredAt,
		IsActive:     true,
	}

	if err := app.DB.Create(&updatedSession).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// revoke old refresh token
	if err := app.DB.Model(&session).Updates(models.Session{IsActive: false}).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Refresh token success", RefreshTokenResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresIn:  AccessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: refreshPayload.ExpiredAt,
	})
}
