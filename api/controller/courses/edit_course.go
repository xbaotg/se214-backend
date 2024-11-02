package courses

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EditCourseRequest struct {
	ID               uuid.UUID  `json:"course_id" binding:"required"`
	CourseTeacherID  uuid.UUID  `json:"course_teacher_id"`
	DepartmentID     uuid.UUID  `json:"course_department"`
	CourseName       string     `json:"course_name"`
	CourseFullname   string     `json:"course_fullname"`
	CourseCredit     int32      `json:"course_credit"`
	CourseYear       int32      `json:"course_year"`
	CourseSemester   int32      `json:"course_semester"`
	CourseStartShift int32      `json:"course_start_shift"`
	CourseEndShift   int32      `json:"course_end_shift"`
	CourseDay        models.Day `json:"course_day"`
	MaxEnroll        int32      `json:"max_enroll"`
	CurrentEnroll    int32      `json:"current_enroll"`
	CourseRoom       string     `json:"course_room"`
}

func EditCourse(c *gin.Context, app *bootstrap.App) {
	// validate request
	req := EditCourseRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	// get current user from session
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// check if user is admin
	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Permission denied", nil)
		return
	}

	// update course with field has value
	course := models.Course{
		ID: req.ID,
	}
	if err := app.DB.First(&course, req.ID).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	reqValue := reflect.ValueOf(req)
	courseValue := reflect.ValueOf(&course).Elem()
	reqType := reqValue.Type()

	for i := 0; i < reqValue.NumField(); i++ {
		field := reqValue.Field(i)
		fieldName := reqType.Field(i).Name

		// Only update if the field is not zero value
		if !field.IsZero() {
			destField := courseValue.FieldByName(fieldName)
			if destField.CanSet() {
				destField.Set(field)
			}
		}
	}

	// update course
	if err := app.DB.Save(&course).Error; err != nil {
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Course updated", course)
}
