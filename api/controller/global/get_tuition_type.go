package global

import (
	"be/bootstrap"
	"be/internal"
	// "be/models"


	"github.com/gin-gonic/gin"
)

// GetTuitionType godoc
// @Summary Get tuition type
// @Description Get tuition type
// @Tags Global
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /global/tuition_type [get]
func GetTuitionType(c *gin.Context, app *bootstrap.App) {
	// sess, _ := c.Get("session")
	// session := sess.(models.Session)

	// // get user info
	// user := models.User{
	// 	ID: session.UserID,
	// }
	// if err := app.DB.First(&user).Error; err != nil {
	// 	internal.Respond(c, 500, false, "Internal server error", nil)
	// 	return
	// }

	internal.Respond(c, 200, true, "Tuition type", gin.H{"type": app.Config.TuitionType, "cost": app.Config.TuitionCost})
}
