package coursesmanagament

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRegisterCourseRequest struct {
	CourseID       uuid.UUID `json:"course_id" binding:"required"`
	CourseYear     int32     `json:"course_year"`
	CourseSemester int32     `json:"course_semester"`
}

func UserRegisterCourse(c *gin.Context, app *bootstrap.App) {
	// validate request
	req := UserRegisterCourseRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	// validate user
	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// validate course
	course := models.Course{ID: req.CourseID}
	if err := app.DB.First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Respond(c, 404, false, "Course not found", nil)
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
	}

	// check prerequisite
	if err := CheckPrerequisite(app, &user, &course); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	// check if user has already registered the course
	// registeredCourse := models.RegisteredCourse{UserID: user.ID, CourseID: course.ID}

	// register course
	registeredCourse := models.RegisteredCourse{
		UserID:         user.ID,
		CourseID:       course.ID,
		CourseYear:     req.CourseYear,
		CourseSemester: req.CourseSemester,
		Status:         models.CoStatusProgressing,
	}
	if err := app.DB.Create(&registeredCourse).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Register course successfully", nil)
}
