package tuition

import (
	"be/internal"
	"be/bootstrap"
	"be/models"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PayTuitionRequest struct {
	TuitionID uuid.UUID `json:"tuition_id"`
	Pay 	 int    `json:"pay"`
	Year 	 int    `json:"year"`
	Semester int    `json:"semester"`
}

// PayTuition godoc
// @Summary Pay tuition
// @Description Pay tuition
// @Tags Tuition
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param PayTuitionRequest body PayTuitionRequest true "PayTuitionRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /tuition/pay [post]
func PayTuition(c *gin.Context, app *bootstrap.App) {
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

	var req PayTuitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if app.State != bootstrap.DONE {
		internal.Respond(c, 403, false, "Chỉ có thể dùng chức năng này khi máy chủ ở trạng thái DONE", nil)
		return
	}

	tuition := models.Tuition{
		ID: req.TuitionID,
		Year: int32(req.Year),
		Semester: int32(req.Semester),
	}

	if err := app.DB.Where(&tuition).First(&tuition).Error; err != nil {
		internal.Respond(c, 400, false, "Không tìm thấy học phí", nil)
		return
	}

	if tuition.TuitionStatus == models.TuStatusPaid {
		internal.Respond(c, 400, false, "Học phí đã được thanh toán", nil)
		return
	}

	tuition.Paid += int32(req.Pay)
	var remaining int32 
	remaining = 0
	if tuition.Paid >= tuition.Tuition {
		remaining = tuition.Paid - tuition.Tuition
		tuition.TuitionStatus = models.TuStatusPaid
		tuition.Paid = tuition.Tuition
	}

	if err := app.DB.Save(&tuition).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Thanh toán học phí thành công", gin.H{
		"tuition": tuition,
		"remaining": remaining,
	})
}