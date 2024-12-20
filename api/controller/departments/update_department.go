package departments

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateDepartmentRequest struct {
	DepartmentName string `json:"department_name" binding:"required"`
	DepartmentCode string `json:"department_code" binding:"required"`
	DepartmentID   uuid.UUID `json:"department_id" binding:"required"`
}

// @Summary Update department
// @Description Update department
// @Tags Department
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param UpdateDepartmentRequest body UpdateDepartmentRequest true "UpdateDepartmentRequest"
// @Router /department/update [put]
func UpdateDepartment(c *gin.Context, app *bootstrap.App) {
	req := UpdateDepartmentRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

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

	department := models.Department{
		ID: req.DepartmentID,
	}

	if err := app.DB.First(&department, req.DepartmentID).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if department.DepartmentCode != req.DepartmentCode {
		if ra := app.DB.Table("departments").Where("department_code = ?", req.DepartmentCode).First(&models.Department{}); ra.RowsAffected > 0 {
			internal.Respond(c, 403, false, "Khoa đã tồn tại, vui lòng kiểm tra lại", nil)
			return
		}
	}

	department.DepartmentName = req.DepartmentName
	department.DepartmentCode = req.DepartmentCode

	if err := app.DB.Transaction(func (tx *gorm.DB) error {
		if err := tx.Save(&department).Error; err != nil {
			return err
		}

		return nil		
	}); err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Cập nhật khoa thành công", nil)
}
