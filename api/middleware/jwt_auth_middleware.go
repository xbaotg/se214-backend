package middleware

import (
	"be/bootstrap"
	"be/model"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SessionMiddleware(app *bootstrap.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the refresh token from the request header
		refreshToken := c.GetHeader("Authorization")
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, model.Response{
				Status:  false,
				Message: "Missing refresh token",
			})
			c.Abort()
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
			c.Abort()
			return
		}

		// Get the session by the token's ID
		session, err := app.DB.GetSessionBySessionId(c, refreshPayload.ID)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, model.Response{
					Status:  false,
					Message: "Session not found",
				})
				c.Abort()
				return
			}
			app.Logger.Error().Err(err).Msg("Failed to get session")
			c.JSON(http.StatusInternalServerError, model.Response{
				Status:  false,
				Message: "Internal server error",
			})
			c.Abort()
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

		if (!session.IsActive) || (session.RefreshToken != refreshToken) {
			c.JSON(401, model.Response{
				Status:  false,
				Message: "Invalid token",
			})
			return
		}

		// Pass on to the next handler
		c.Set("session", session)
		c.Next()
	}
}
