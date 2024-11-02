package middleware

import (
	"be/bootstrap"
	"be/internal"
	"be/model"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SessionMiddleware(app *bootstrap.App, shouldCancel bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the refresh token from the request header
		refreshToken := c.GetHeader("Authorization")
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, model.Response{
				Status:  false,
				Message: "Missing refresh token",
			})

			if shouldCancel {
				c.Abort()
			}

			return
		}

		// Verify the refresh token
		refreshToken = refreshToken[len("Bearer "):]
		refreshPayload, err := app.TokenMaker.VerifyToken(refreshToken)
		if err != nil {
			app.Logger.Error().Msg(err.Error())

			c.JSON(http.StatusUnauthorized, model.Response{
				Status:  false,
				Message: "Invalid token",
			})
			if shouldCancel {
				c.Abort()
			}
			return
		}

		// Get the session by the token's ID
		session, err := app.DB.GetSessionBySessionId(c, refreshPayload.ID)

		if err != nil {
			if err == sql.ErrNoRows {
				internal.Respond(c, 404, false, "Session not found", nil)
				if shouldCancel {
					c.Abort()
				}
				return
			}

			app.Logger.Error().Err(err).Msg("Failed to get session")
			internal.Respond(c, 500, false, "Internal server error", nil)
			if shouldCancel {
				c.Abort()
			}
			return
		}

		// check if refresh token is valid
		if time.Now().After(session.ExpiresIn) {
			internal.Respond(c, 401, false, "Refresh token expired", nil)
			if shouldCancel {
				c.Abort()
			}
			return
		}

		if (!session.IsActive) || (session.RefreshToken != refreshToken) {
			internal.Respond(c, 401, false, "Invalid token", nil)
			if shouldCancel {
				c.Abort()
			}
			return
		}

		// Pass on to the next handler
		c.Set("session", session)
		c.Next()
	}
}
