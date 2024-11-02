package course

import (
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateCourseRequest struct {
	CourseTeacherID  uuid.UUID `json:"course_teacher_id" binding:"required"`
	CourseDepartment uuid.UUID `json:"course_department" binding:"required"`
	CourseName       string    `json:"course_name" binding:"required,min=3,max=100"`
	CourseFullname   string    `json:"course_fullname" binding:"required,min=5,max=200"`
	CourseCredit     int32     `json:"course_credit" binding:"required,min=1,max=12"`
	CourseYear       int32     `json:"course_year" binding:"required,min=2000,max=2100"`
	CourseSemester   int32     `json:"course_semester" binding:"required,oneof=1 2 3"`
	CourseStartShift int32     `json:"course_start_shift" binding:"required,min=1,max=10"`
	CourseEndShift   int32     `json:"course_end_shift" binding:"required,min=1,max=10"`
	CourseDay        sqlc.Day  `json:"course_day" binding:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"`
	MaxEnroll        int32     `json:"max_enroll" binding:"required,min=1,max=1000"`
	CurrentEnroll    int32     `json:"current_enroll" binding:"required,min=0,max=1000"`
	CourseRoom       string    `json:"course_room" binding:"required,min=2,max=50"`
}

// Create course
// @Summary Create course
// @Description Create course
// @Tags Course
// @Accept json
// @Produce json
// @Param CreateCourseRequest body CreateCourseRequest true "CreateCourseRequest"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 403 {object} model.Response
// @Failure 500 {object} model.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /course/create [post]
func CreateCourse(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(sqlc.Session)

	// request validation
	req := CreateCourseRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
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
	user, err := app.DB.GetUserById(c, session.UserID)
	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	// check user role
	if user.UserRole != sqlc.RoleAdmin {
		internal.Respond(c, 403, false, "Permission denied", nil)
		return
	}

	// create course
	course, err := app.DB.CreateCourse(c, sqlc.CreateCourseParams{
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
		MaxEnroller:      req.MaxEnroll,
		CurrentEnroller:  req.CurrentEnroll,
		CourseRoom:       req.CourseRoom,
		CreatedAt:        internal.GetCurrentTime(),
		UpdatedAt:        internal.GetCurrentTime(),
	})

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "Course created", course)
}
