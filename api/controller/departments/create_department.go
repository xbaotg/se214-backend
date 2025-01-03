package departments

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

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
// @Success 200 {object} models.Department
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /department/create [post]
func CreateDepartment(c *gin.Context, app *bootstrap.App) {
	// validate request
	req := CreateDepartmentRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, "Hãy điền đầy đủ thông tin", nil)
		return
	}

	// get current user from session
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	// check if user is admin
	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	// department, err := app.DB.CreateDepartment(c, sqlc.CreateDepartmentParams{
	// 	DepartmentName: req.DepartmentName,
	// 	DepartmentCode: req.DepartmentCode,
	// 	CreatedAt:      internal.GetCurrentTime(),
	// 	UpdatedAt:      internal.GetCurrentTime(),
	// })

	departmentToCreate := models.Department{
		ID:             internal.GenerateUUID(),
		DepartmentName: req.DepartmentName,
		DepartmentCode: req.DepartmentCode,
	}

	if ra := app.DB.Table("departments").Where("department_code = ?", req.DepartmentCode).First(&models.Department{}); ra.RowsAffected > 0 {
		internal.Respond(c, 403, false, "Khoa đã tồn tại, vui lòng kiểm tra lại", nil)
		return
	}

	if err := app.DB.Create(&departmentToCreate).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	departmentReponse := DepartmentResponse{
		DepartmentID:   departmentToCreate.ID.String(),
		DepartmentName: departmentToCreate.DepartmentName,
		DepartmentCode: departmentToCreate.DepartmentCode,
		CreatedAt:      departmentToCreate.CreatedAt.String(),
		UpdatedAt:      departmentToCreate.UpdatedAt.String(),
	}

	internal.Respond(c, 200, true, "Tạo khoa thành công", departmentReponse)
}
