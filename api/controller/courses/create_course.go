package courses

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"be/api/schemas"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"slices"
	"fmt"
)

type CreateCourseRequest struct {
	// CourseCode       string     `json:"course_code"`
	CourseTeacherID  uuid.UUID  `json:"course_teacher_id" binding:"required"`
	CourseDepartment uuid.UUID  `json:"course_department" binding:"required"`
	CourseName       string     `json:"course_name" binding:"required,min=3,max=100"`
	CourseFullname   string     `json:"course_fullname" binding:"required,min=5,max=200"`
	CourseCredit     int32      `json:"course_credit" binding:"required,min=1,max=12"`
	CourseYear       int32      `json:"course_year" binding:"required,min=2000,max=2100"`
	CourseSemester   int32      `json:"course_semester" binding:"required,oneof=1 2 3"`
	CourseStartShift int32      `json:"course_start_shift" binding:"required,min=1,max=10"`
	CourseEndShift   int32      `json:"course_end_shift" binding:"required,min=1,max=10"`
	CourseDay        models.Day `json:"course_day" binding:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	MaxEnroll        int32      `json:"max_enroll" binding:"required,min=1,max=1000"`
	CurrentEnroll    int32      `json:"current_enroll" binding:"min=0,max=1000"`
	CourseRoom       string     `json:"course_room" binding:"required,min=2,max=50"`
}


// Create course
// @Summary Create course
// @Description Create course
// @Tags Course
// @Accept json
// @Produce json
// @Param CreateCourseRequest body CreateCourseRequest true "CreateCourseRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /course/create [post]
func CreateCourse(c *gin.Context, app *bootstrap.App) {
	if app.State != bootstrap.SETUP {
		internal.Respond(c, 403, false, fmt.Sprintf("Server is not in setup state, current state is %s", app.State), nil)
		return
	}
	sess, _ := c.Get("session")
	session := sess.(models.Session)

	// request validation
	req := CreateCourseRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// internal.Respond(c, 400, false, "Please fill all required fields", nil)
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	if req.CurrentEnroll > req.MaxEnroll {
		internal.Respond(c, 400, false, "Current enroll must be less than max enroll", nil)
		return
	}

	if req.CourseStartShift >= req.CourseEndShift {
		internal.Respond(c, 400, false, "Course start shift must be less than course end shift", nil)
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

	// check user role
	if !slices.Contains([]models.Role{models.RoleAdmin, models.RoleLecturer}, user.UserRole) {
		internal.Respond(c, 403, false, "Permission denied", nil)
		return
	}

	// currentCourses := []models.Course{
	// 	models.Course{
	// 		CourseDay:      req.CourseDay,
	// 		CourseYear:     req.CourseYear,
	// 		CourseSemester: req.CourseSemester,
	// 		CourseRoom:     req.CourseRoom,
	// 	},
	// }

	// if err := app.DB.Where(currentCourses).Find(&currentCourses).Error; err != nil {
	// 	app.Logger.Error().Err(err).Msg(err.Error())
	// 	internal.Respond(c, 500, false, "Internal server error", nil)
	// 	return
	// }

	// if len(currentCourses) > 0 {
	// 	for _, course := range currentCourses {
	// 		for i := course.CourseStartShift; i <= course.CourseEndShift; i++ {
	// 			if i >= req.CourseStartShift && i <= req.CourseEndShift {
	// 				internal.Respond(c, 400, false, "Course shift is not available", nil)
	// 				return
	// 			}
	// 		}
	// 	}
	// }

	course := models.Course{
		ID:               internal.GenerateUUID(),
		CourseTeacherID:  req.CourseTeacherID,
		DepartmentID:     req.CourseDepartment,
		CourseName:       req.CourseName,
		CourseFullname:   req.CourseFullname,
		CourseCredit:     req.CourseCredit,
		CourseYear:       req.CourseYear,
		CourseSemester:   req.CourseSemester,
		CourseStartShift: req.CourseStartShift,
		CourseEndShift:   req.CourseEndShift,
		CourseDay:        req.CourseDay,
		Confirmed:        false,
		MaxEnroller:      req.MaxEnroll,
		CurrentEnroller:  req.CurrentEnroll,
		CourseRoom:       req.CourseRoom,
	}
	if user.UserRole == models.RoleAdmin {
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
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}

		if len(currentCourses) > 0 {
			for _, course := range currentCourses {
				for i := course.CourseStartShift; i <= course.CourseEndShift; i++ {
					if i >= req.CourseStartShift && i <= req.CourseEndShift {
						internal.Respond(c, 400, false, "Course shift is not available", nil)
						return
					}
				}
			}
		}
		course.Confirmed = true
	}

	if err := app.DB.Create(&course).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	courseResponse := schemas.CreateCourseResponse{
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
	}

	internal.Respond(c, 200, true, "Course created", courseResponse)
}
