package tuition

import (
	"be/internal"
	"be/bootstrap"
	"be/models"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateTuitionRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Year 	 int32    `json:"year" binding:"required"`
	Semester int32    `json:"semester" binding:"required"`
	Tuition  int32    `json:"tuition"`
	Deadline string `json:"deadline"`
	TuitionStatus models.TuStatus `json:"tuition_status"`
}

// UpdateTuition godoc
// @Summary Update tuition
// @Description Update tuition
// @Tags Tuition
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param UpdateTuitionRequest body UpdateTuitionRequest true "UpdateTuitionRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /tuition/update_tuition [put]
func UpdateTuition(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

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

	var req UpdateTuitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, "Hãy điền đầy đủ thông tin", nil)
		return
	}

	deadline, err := time.Parse("2006-01-02T15:04:05.000Z", req.Deadline)
	if err != nil {
		internal.Respond(c, 400, false, "Deadline không hợp lệ", nil)
		return
	}

	tuition := models.Tuition{
		UserID: req.UserID,
		Year: req.Year,
		Semester: req.Semester,
	}

	if err := app.DB.Where(&tuition).First(&tuition).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	tuition.Tuition = req.Tuition
	tuition.TuitionDeadline = deadline
	tuition.TuitionStatus = req.TuitionStatus

	if err := app.DB.Save(&tuition).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, http.StatusOK, true, "Cập nhật học phí thành công", tuition)
}