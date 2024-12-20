package global

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"strconv"


	"github.com/gin-gonic/gin"
)

type SetTuitionTypeRequest struct {
	Type string `json:"type" binding:"required"`
	Cost int `json:"cost" binding:"required"`
}

// SetTuitionType godoc
// @Summary Set tuition type
// @Description Set tuition type
// @Tags Global
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param SetTuitionTypeRequest body SetTuitionTypeRequest true "SetTuitionTypeRequest"
// @Router /global/tuition_type [post]
func SetTuitionType(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	var req SetTuitionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, "Hãy điền đầy đủ thông tin", nil)
		return
	}

	// get user info
	user := models.User{
		ID: session.UserID,
	}
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}


	app.Config.TuitionType = req.Type
	app.Config.TuitionCost = req.Cost
	UpdateEnvFile("TUITION_TYPE", req.Type)
	UpdateEnvFile("TUITION_COST", strconv.Itoa(req.Cost))
	internal.Respond(c, 200, true, "Cập nhật loại học phí thành công", nil)
}
