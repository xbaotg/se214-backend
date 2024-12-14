package tuition

import (
	"be/internal"
	"be/bootstrap"
	"be/models"
	"be/api/schemas"
	// "net/http"
	// "slices"
	
	"github.com/gin-gonic/gin"
)

type CalTuitionRequest struct {
	Year int `json:"year" binding:"required"`
	Semester int `json:"semester" binding:"required"`
}

// CalTuition godoc
// @Summary Calculate tuition
// @Description Calculate tuition
// @Tags Tuition
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param CalTuitionRequest body CalTuitionRequest true "CalTuitionRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /tuition/cal_tuition [post]
func CalTuition(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	// if app.State == bootstrap.ACTIVE {
	// 	internal.Respond(c, 403, false, "Cannot calculate tuition in this state", nil)
	// 	return
	// }
	var req CalTuitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, "Invalid request", nil)
		return
	}

	// get user info
	user := models.User{
		ID: session.UserID,
	}
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	var courses []models.Course
	var credit int32 = 0
	if err := app.DB.Select("courses.*").Table(
		"registered_courses").Joins("JOIN courses ON courses.id = registered_courses.course_id").Where(
		"registered_courses.user_id = ? AND registered_courses.course_year = ? AND registered_courses.course_semester = ?",
		user.ID, req.Year, req.Semester).Find(&courses).Error; err != nil {
		

		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}	

	coursesResponse := []schemas.CreateCourseResponse{}
	if len(courses) == 0 {
		internal.Respond(c, 200, true, "User has not registered any course", gin.H{"tuition": -1})
		return
	} else {
		for _, course := range courses {
			credit += course.CourseCredit
			coursesResponse = append(coursesResponse, schemas.CreateCourseResponse{
			ID:               course.ID,
			CourseTeacherID:  course.CourseTeacherID,
			CourseDepartment: course.DepartmentID,
			CourseName:       course.CourseName,
			CourseFullname:   course.CourseFullname,
			CourseCredit:     course.CourseCredit,
			CourseYear:       course.CourseYear,
			CourseSemester:   course.CourseSemester,
			CourseStartShift: course.CourseStartShift,
			CourseEndShift:   course.CourseEndShift,
			CourseDay:        course.CourseDay,
			MaxEnroll:        course.MaxEnroller,
			CurrentEnroll:    course.CurrentEnroller,
			CourseRoom:       course.CourseRoom,
			})
		}
	}
	if app.Config.TuitionType == "buffet" {
		internal.Respond(c, 200, true, "Tuition", gin.H{
			"tuition": app.Config.TuitionCost,
			"courses": coursesResponse,
			"credit": credit,
		})
	} else {		
		internal.Respond(c, 200, true, "Tuition", gin.H{
			"tuition": credit * int32(app.Config.TuitionCost),
			"courses": coursesResponse,
			"credit": credit,
		})
	}
}