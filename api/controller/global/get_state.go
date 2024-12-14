package global

import (
	"be/internal"
	"be/bootstrap"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetState godoc
// @Summary Get the state of the app
// @Description Get the state of the app
// @Tags Global
// @Produce  json
// @Success 200 {object} string
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /state [get]
func GetState(c *gin.Context, app *bootstrap.App) {
	internal.Respond(c, http.StatusOK, true, "State", gin.H{"state": app.State})
}