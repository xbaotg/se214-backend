package departments

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"

	"github.com/gin-gonic/gin"
)

// @Summary List department
// @Description List department
// @Tags Department
// @Accept json
// @Produce json
// @Success 200 {object} sqlc.Department
// @Failure 403 {object} model.Response
// @Failure 500 {object} model.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /department/list [get]
func ListDepartment(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(sqlc.Session)
	user, err := app.DB.GetUserById(c, session.UserID)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	if user.UserRole != sqlc.RoleAdmin {
		internal.Respond(c, 403, false, "Forbidden", nil)
		return
	}

	departments, err := app.DB.ListDepartments(c)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "List department success", departments)
}
