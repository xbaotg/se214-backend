package departments

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

type DepartmentResponse struct {
	DepartmentID   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
	DepartmentCode string `json:"department_code"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

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
	departments := []models.Department{}
	if err := app.DB.Find(&departments).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// map department to response
	departmentResponses := []DepartmentResponse{}
	for _, department := range departments {
		departmentResponses = append(departmentResponses, DepartmentResponse{
			DepartmentID:   department.ID.String(),
			DepartmentName: department.DepartmentName,

			DepartmentCode: department.DepartmentCode,
			CreatedAt:      department.CreatedAt.String(),
			UpdatedAt:      department.UpdatedAt.String(),
		})
	}

	internal.Respond(c, 200, true, "List department success", departmentResponses)
}
