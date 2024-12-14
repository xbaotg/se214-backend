package coursesmanagament

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
)

type GetCreditRequest struct {
	Year int `json:"year" binding:"required"`
	Semester int `json:"semester" binding:"required"`
}

// GetCredit godoc
// @Summary Get credit
// @Description Get credit
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param GetCreditRequest body GetCreditRequest true "GetCreditRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/credit [post]
func GetCredit(c *gin.Context, app *bootstrap.App) {
	req := GetCreditRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, "Invalid request", nil)
		return
	}

	sess, _ := c.Get("session")
	session := sess.(models.Session)

	user := models.User{
		ID: session.UserID,
	}

	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	var credit int64
	if err := app.DB.Select("SUM(courses.course_credit) as credit").Table(
		"registered_courses").Joins("JOIN courses ON courses.id = registered_courses.course_id").Where(
		"registered_courses.user_id = ? AND registered_courses.course_year = ? AND registered_courses.course_semester = ?",
		user.ID, req.Year, req.Semester).Scan(&credit).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Credit", gin.H{"credit": credit})
}
