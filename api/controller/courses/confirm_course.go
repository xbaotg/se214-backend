package courses

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConfirmCourse docs
// @Summary Confirm Course
// @Description Confirm course
// @Tags Courses
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param course_id query string true "Course ID"
// @Router /course/confirm [put]
func ConfirmCourse(c *gin.Context, app *bootstrap.App) {
	if app.State != bootstrap.SETUP {
		internal.Respond(c, 403, false, fmt.Sprintf("Máy chủ không ở trạng thái SETUP, trạng thái hiện tại là %s", app.State), nil)
		return
	}
	// get current user from session
	courseID := c.Query("course_id")
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}
	courseIDUUID, err := uuid.Parse(courseID)
	if err != nil {
		internal.Respond(c, 400, false, "Mã khóa học không hợp lệ", nil)
		return
	}

	// check if user is admin
	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	// update course with field has value
	course := models.Course{
		ID: courseIDUUID,
	}
	if err := app.DB.First(&course).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}
	currentCourses := []models.Course{
		models.Course{
			CourseDay:      course.CourseDay,
			CourseYear:     course.CourseYear,
			CourseSemester: course.CourseSemester,
			CourseRoom:     course.CourseRoom,
		},
	}

	if err := app.DB.Where(currentCourses).Where("confirmed = ?", true).Find(&currentCourses).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	app.Logger.Info().Msgf("currentCourses: %v", currentCourses)

	if len(currentCourses) > 0 {
		for _, c_ := range currentCourses {
			for i := c_.CourseStartShift; i <= c_.CourseEndShift; i++ {
				if i >= course.CourseStartShift && i <= course.CourseEndShift {
					internal.Respond(c, 400, false, "Ca học bị trùng", nil)
					return
				}
			}
		}
	}

	if err := app.DB.Model(&course).Update("confirmed", true).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Xác nhận khóa học thành công", course)
}