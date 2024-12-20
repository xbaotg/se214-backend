package coursesmanagament

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sync"
)

type UserRegisterCourseRequest struct {
	CourseID       uuid.UUID `json:"course_id" binding:"required"`
	CourseYear     int32     `json:"course_year"`
	CourseSemester int32     `json:"course_semester"`
}

var (
	ErrAlreadyRegistered = errors.New("Đã đăng ký khóa học")
	ErrCourseFull        = errors.New("Khóa học đã đủ số lượng học viên")
)

// courseid chan string to lock user register course
var userLocks = make(map[uuid.UUID]chan uuid.UUID)
var lockMutex sync.RWMutex

// UserRegisterCourse docs
// @Summary User register course
// @Description User register course
// @Tags User
// @Accept json
// @Produce json
// @Param UserRegisterCourseRequest body UserRegisterCourseRequest true "UserRegisterCourseRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /user/course/register [post]
func UserRegisterCourse(c *gin.Context, app *bootstrap.App) {
	if app.State != bootstrap.ACTIVE {
		internal.Respond(c, 403, false, fmt.Sprintf("Máy chủ không ở trạng thái ACTIVE, trạng thái hiện tại là %s", app.State), nil)
		return
	}
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
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	// validate course
	course := models.Course{ID: req.CourseID}
	if err := app.DB.First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			internal.Respond(c, 404, false, "Không tìm thấy khóa học", nil)
			return
		}

		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}

	// check prerequisite
	if co, err := CheckPrerequisite(app, &user, &course); err != nil {
		internal.Respond(c, 406, false, err.Error(), gin.H{"course": co})
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

	lockMutex.Lock()
    if _, ok := userLocks[user.ID]; !ok {
        userLocks[user.ID] = make(chan uuid.UUID, 1) // Buffer size 1
    }
    lockMutex.Unlock()

    // Block and wait until previous registration is done
    userLocks[user.ID] <- course.ID // This will block if channel is full
    
    // Ensure we release the lock when we're done
    defer func() {
        <-userLocks[user.ID] // Remove our registration from the channel
    }()
	// check if already have a course in the same year, semester, day, shift
	var courses []models.Course	
	if err := app.DB.Table("courses").Select("courses.*",
		).Joins(
		"JOIN registered_courses ON courses.id = registered_courses.course_id",
		).Where("registered_courses.user_id = ? AND courses.course_day = ? AND registered_courses.course_year = ? AND registered_courses.course_semester = ?", 
		user.ID, course.CourseDay, req.CourseYear, req.CourseSemester).Find(&courses).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if len(courses) > 0 {
		for _, course_ := range courses {
			for i := course_.CourseStartShift; i <= course_.CourseEndShift; i++ {
				if i >= course.CourseStartShift && i <= course.CourseEndShift {
					internal.Respond(c, 400, false, "Ca học bị trùng", nil)
					return
				}
			}
		}
	}	

	// Update current enroller of the course
	if err := app.DB.Transaction(func(tx *gorm.DB) error {

		course.CurrentEnroller++
		if course.CurrentEnroller > course.MaxEnroller {
			return ErrCourseFull
		}
		if err := tx.Create(&registeredCourse).Error; err != nil {
			return ErrAlreadyRegistered
		}

		if err := tx.Save(&course).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		switch err {
		case ErrAlreadyRegistered:
			internal.Respond(c, 400, false, ErrAlreadyRegistered.Error(), nil)
			return
		case ErrCourseFull:
			internal.Respond(c, 400, false, ErrCourseFull.Error(), nil)
			return
		default:
			internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
			return
		}
	}

	internal.Respond(c, 200, true, "Đăng ký khóa học thành công", registeredCourse)
}
