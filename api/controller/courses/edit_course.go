package courses

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"reflect"
	"fmt"

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

// @Summary Edit course
// @Description Edit course
// @Tags Course
// @Accept json
// @Produce json
// @Param EditCourseRequest body EditCourseRequest true "EditCourseRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /course/edit [put]
func EditCourse(c *gin.Context, app *bootstrap.App) {
	if app.State != bootstrap.SETUP {
		internal.Respond(c, 403, false, fmt.Sprintf("Máy chủ không ở trạng thái SETUP, trạng thái hiện tại là %s", app.State), nil)
		return
	}

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
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	// check if user is admin
	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	// update course with field has value
	course := models.Course{
		ID: req.ID,
	}
	if err := app.DB.First(&course, req.ID).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}


	courseName := models.AllCourses{CourseName: req.CourseName, Status: true}
	if err := app.DB.Table("all_courses").Where(courseName).FirstOrCreate(&courseName).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	currentCourses := []models.Course{
		models.Course{
			CourseDay:      req.CourseDay,
			CourseYear:     req.CourseYear,
			CourseSemester: req.CourseSemester,
			CourseRoom:     req.CourseRoom,
		},
	}

	if err := app.DB.Where(currentCourses).Find(&currentCourses).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if len(currentCourses) > 0 {
		for _, courseLoop := range currentCourses {
			for i := courseLoop.CourseStartShift; i <= courseLoop.CourseEndShift; i++ {
				if i >= req.CourseStartShift && i <= req.CourseEndShift && req.ID != courseLoop.ID {
					internal.Respond(c, 400, false, "Ca học bị trùng", nil)
					return
				}
			}
		}
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
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Cập nhật khóa học thành công", nil)
}
