package middleware

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SessionMiddleware(app *bootstrap.App, shouldCancel bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Get the refresh token from the request header
		refreshToken := c.GetHeader("Authorization")
		if refreshToken == "" {
			if shouldCancel {
				internal.Respond(c, 400, false, "Missing refresh token", nil)
				c.Abort()
			}

			return
		}

		// Verify the refresh token
		refreshToken = refreshToken[len("Bearer "):]
		refreshPayload, err := app.TokenMaker.VerifyToken(refreshToken)
		if err != nil {
			app.Logger.Error().Msg(err.Error())

			if shouldCancel {
				internal.Respond(c, 401, false, "Invalid token", nil)
				c.Abort()
			}
			return
		}

		// Get the session by the token's ID
		// session, err := app.DB.GetSessionBySessionId(c, refreshPayload.ID)
		session := models.Session{ID: refreshPayload.ID}
		if err := app.DB.First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if shouldCancel {
					internal.Respond(c, 404, false, "Session not found", nil)
					c.Abort()
				}
				return
			}

			app.Logger.Error().Err(err).Msg("Failed to get session")
			if shouldCancel {
				internal.Respond(c, 500, false, "Internal server error", nil)
				c.Abort()
			}
			return
		}

		// check if refresh token is valid
		if time.Now().After(session.ExpiresIn) {
			if shouldCancel {
				internal.Respond(c, 401, false, "Refresh token expired", nil)
				c.Abort()
			}
			return
		}

		if (!session.IsActive) || (session.RefreshToken != refreshToken) {
			if shouldCancel {
				internal.Respond(c, 401, false, "Invalid token", nil)
				c.Abort()
			}
			return
		}

		// Pass on to the next handler
		c.Set("session", session)
		c.Next()
	}
}
