package lecturer

import (
	"be/api/schemas"
	"be/bootstrap"
	"be/internal"
	"be/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// List lecturer register courses
// @Summary List lecturer register courses
// @Description List lecturer register courses
// @Tags Lecturer
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /lecturer/course/register/list [get]
func ListLecturerRegisterCourses(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	// if user.Role != models.Lecturer {
	// 	internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
	// 	return
	// }

	courses := []models.Course{}
	if user.UserRole == models.RoleLecturer {
		if err := app.DB.Where("course_teacher_id = ? and confirmed = ?", user.ID, false).Find(&courses).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				internal.Respond(c, 404, false, "Không tìm thấy khóa học", nil)
				return
			}
			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
			return
		}
	} else if user.UserRole == models.RoleAdmin {
		if err := app.DB.Where("confirmed = ?", false).Find(&courses).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				internal.Respond(c, 404, false, "Không tìm thấy khóa học", nil)
				return
			}
			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
			return
		}
	} else {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	var courseResponse []schemas.CreateCourseResponse
	for _, course := range courses {
		courseResponse = append(courseResponse, schemas.CreateCourseResponse{
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
	if len(courseResponse) == 0 {
		internal.Respond(c, 200, true, "Không có khóa học nào", []schemas.CreateCourseResponse{})
		return
	}
	internal.Respond(c, 200, true, "Lấy danh sách khóa học thành công", courseResponse)
}