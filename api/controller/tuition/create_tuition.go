package tuition

import (
	"be/internal"
	"be/bootstrap"
	"be/models"
	"net/http"
	"time"
	"sync"
	
	"github.com/gin-gonic/gin"
)

const MAX_GOROUTINES = 10

type CreateTuitionRequest struct {
	Year int `json:"year" binding:"required"`
	Semester int `json:"semester" binding:"required"`
	Deadline string `json:"deadline" binding:"required"`
}

// CreateTuition godoc
// @Summary Create tuition
// @Description Create tuition
// @Tags Tuition
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param CreateTuitionRequest body CreateTuitionRequest true "CreateTuitionRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /tuition/create_tuition [post]
func CreateTuition(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	// get user info
	user := models.User{
		ID: session.UserID,
	}
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Permission denied", nil)
		return
	}

	var req CreateTuitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
// 2025-05-13T17:00:00.000Z"
	deadline, err := time.Parse("2006-01-02T15:04:05.000Z", req.Deadline)
	if err != nil {
		internal.Respond(c, 400, false, "Invalid deadline", nil)
		return
	}

	if app.State != bootstrap.DONE {
		internal.Respond(c, 403, false, "Can only create tuition when server is in DONE state", nil)
		return
	}

	var users []models.User
	if err := app.DB.Table("users").Where("user_role = ?", models.RoleUser).Find(&users).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	var wg sync.WaitGroup
	lenWg := 0

	for _, u := range users {
		wg.Add(1)
		lenWg +=1
		go func() {
			defer wg.Done()
			defer func() {
				lenWg -= 1
			}()
			var courses []models.Course
			if err := app.DB.Select("courses.*").Table(
				"registered_courses").Joins("JOIN courses ON courses.id = registered_courses.course_id").Where(
				"registered_courses.user_id = ? AND registered_courses.course_year = ? AND registered_courses.course_semester = ?",
				u.ID, req.Year, req.Semester).Find(&courses).Error; err != nil {
				app.Logger.Error().Err(err).Msg(err.Error())
				return
			}				

			var tuition int
			var credit int
			for _, course := range courses {
				credit += int(course.CourseCredit)
			}
 
			
			if len(courses) > 0 {
				if app.Config.TuitionType == "buffet" {
					tuition = app.Config.TuitionCost
				} else {
					tuition = credit * app.Config.TuitionCost
				}
			} else {
				tuition = -1
			}

			tuitionRecord := models.Tuition{
				UserID: u.ID,
				Year: int32(req.Year),
				Semester: int32(req.Semester),
			}

			if err := app.DB.Where(&tuitionRecord).First(&tuitionRecord).Error; err == nil {
				app.Logger.Info().Msgf("Tuition for user %s in year %d, semester %d already exists", u.ID, req.Year, req.Semester)
				app.Logger.Info().Msgf("%v", tuitionRecord)
				tuitionRecord.Tuition = int32(tuition)
				tuitionRecord.TotalCredit = int32(credit)
				tuitionRecord.TuitionDeadline = deadline
				tuitionRecord.TuitionStatus = models.TuStatusUnpaid

				if err := app.DB.Save(&tuitionRecord).Error; err != nil {
					app.Logger.Error().Err(err).Msg(err.Error())
					return
				}

			} else {
				app.Logger.Info().Msg("Creating new tuition")
				tuitionRecord := models.Tuition{
					ID: internal.GenerateUUID(),
					UserID: u.ID,
					Year: int32(req.Year),
					Semester: int32(req.Semester),
					Tuition: int32(tuition),
					TotalCredit: int32(credit),
					TuitionDeadline: deadline,
					TuitionStatus: models.TuStatusUnpaid,
				}

				if err := app.DB.Create(&tuitionRecord).Error; err != nil {
					app.Logger.Error().Err(err).Msg(err.Error())
					return
				}
			}
		}()
		if lenWg >= MAX_GOROUTINES {
			app.Logger.Info().Msgf("Waiting for %d goroutines to finish", lenWg)
			wg.Wait()
		}
	}

	wg.Wait()

	internal.Respond(c, 200, true, "All tuition created", nil)
}