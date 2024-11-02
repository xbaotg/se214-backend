package departments

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

// @Summary List department
// @Description List department
// @Tags Department
// @Accept json
// @Produce json
// @Success 200 {object} models.Department
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /department/list [get]
func ListDepartment(c *gin.Context, app *bootstrap.App) {
	// sess, _ := c.Get("session")
	// session := sess.(models.Session)
	// user := models.User{ID: session.UserID}

	// if err := app.DB.First(&user).Error; err != nil {
	// 	app.Logger.Error().Err(err).Msg(err.Error())
	// 	internal.Respond(c, 500, false, "Internal server error", nil)
	// 	return
	// }

	// if user.UserRole != models.RoleAdmin {
	// 	internal.Respond(c, 403, false, "Forbidden", nil)
	// 	return
	// }

	departments := []models.Department{}
	if err := app.DB.Find(&departments).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "List department success", departments)
}
