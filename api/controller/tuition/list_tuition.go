package tuition

import (
	"be/internal"
	"be/bootstrap"
	"be/models"
	
	"github.com/gin-gonic/gin"
)

// ListTuition godoc
// @Summary List tuition
// @Description List tuition
// @Tags Tuition
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /tuition/list [get]
func ListTuition(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	// get user info
	user := models.User{
		ID: session.UserID,
	}
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	var tuitions []models.Tuition
	if user.UserRole == models.RoleAdmin {
		if err := app.DB.Find(&tuitions).Error; err != nil {
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}
		internal.Respond(c, 200, true, "Tuition list", tuitions)
		return
	} else {
		if err := app.DB.Where("user_id = ?", user.ID).Find(&tuitions).Error; err != nil {
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}
		internal.Respond(c, 200, true, "Tuition list", tuitions)
		return
	}
}