package departments

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"

	"github.com/gin-gonic/gin"
)

type CreateDepartmentRequest struct {
	DepartmentName string `json:"department_name" binding:"required,min=3,max=100"`
	DepartmentCode string `json:"department_code" binding:"required,alphanum,min=2,max=10"`
}

// @Summary Create department
// @Description Create department
// @Tags Department
// @Accept json
// @Produce json
// @Param CreateDepartmentRequest body CreateDepartmentRequest true "CreateDepartmentRequest"
// @Success 200 {object} sqlc.Department
// @Failure 400 {object} model.Response
// @Failure 403 {object} model.Response
// @Failure 500 {object} model.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /department/create [post]
func CreateDepartment(c *gin.Context, app *bootstrap.App) {
	// validate request
	req := CreateDepartmentRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	// get current user from session
	sess, _ := c.Get("session")
	session := sess.(sqlc.Session)
	user, err := app.DB.GetUserById(c, session.UserID)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// check if user is admin
	if user.UserRole != sqlc.RoleAdmin {
		internal.Respond(c, 403, false, "Forbidden", nil)
		return
	}

	department, err := app.DB.CreateDepartment(c, sqlc.CreateDepartmentParams{
		DepartmentName: req.DepartmentName,
		DepartmentCode: req.DepartmentCode,
		CreatedAt:      internal.GetCurrentTime(),
		UpdatedAt:      internal.GetCurrentTime(),
	})

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Department created", department)
}
