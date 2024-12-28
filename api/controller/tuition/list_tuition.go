package tuition

import (
	"be/internal"
	"be/bootstrap"
	"be/models"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TuitionResponse struct {
	ID              uuid.UUID 
	UserID          uuid.UUID
	Tuition         int32
	Paid            int32
	TotalCredit     int32
	Year            int32
	Semester        int32
	TuitionStatus   models.TuStatus
	TuitionDeadline time.Time
	CreatedAt       time.Time 
	UpdatedAt       time.Time 
	Username 	  string
}

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
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	var tuitions []TuitionResponse
	if user.UserRole == models.RoleAdmin {
		if err := app.DB.Table("tuitions").Select("tuitions.*, users.username").Joins(
			"left join users on tuitions.user_id = users.id",
			).Find(&tuitions).Error; err != nil {
			internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
			return
		}
		internal.Respond(c, 200, true, "Danh sách học phí", tuitions)
	} else {
		if err := app.DB.Table("tuitions").Select("tuitions.*, users.username").Joins(
			"left join users on tuitions.user_id = users.id",
			).Where("tuitions.user_id = ?", session.UserID).Find(&tuitions).Error; err != nil {
			internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
			return
		}
		internal.Respond(c, 200, true, "Danhs sách học phí", tuitions)
	}
	return
}