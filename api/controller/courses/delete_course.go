package courses

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

// type DeleteCourseRequest struct {
// 	CourseID uuid.UUID `json:"course_id" binding:"required"`
// }

// @Summary Delete course
// @Description Delete course
// @Tags Courses
// @Accept json
// @Produce json
// @Param course_id body string true "Course ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /course/delete [delete]
func DeleteCourse(c *gin.Context, app *bootstrap.App) {
	// validate request
	// req := DeleteCourseRequest{}
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	internal.Respond(c, 400, false, err.Error(), nil)
	// 	return
	// }
	if app.State != bootstrap.SETUP {
		internal.Respond(c, 403, false, fmt.Sprintf("Máy chủ không ở trạng thái SETUP, trạng thái hiện tại là %s", app.State), nil)
		return
	}

	courseID_ := c.Param("course_id")
	courseID, err := uuid.Parse(courseID_)
	if err != nil {
		internal.Respond(c, 400, false, "Mã khóa học không hợp lệ", nil)
		return
	}

	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	// check if user exists
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, err.Error(), nil)
		return
	}

	course := models.Course{ID: courseID}
	if err := app.DB.First(&course).Error; err != nil {
		internal.Respond(c, 404, false, "Không tìm thấy khóa học", nil)
		return
	}

	if user.UserRole == models.RoleLecturer {
		if course.CourseTeacherID != user.ID {
			internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
			return
		}

	} else if user.UserRole == models.RoleUser {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	} 

	// delete course
	if err := app.DB.Delete(&course).Error; err != nil {
		internal.Respond(c, 500, false, err.Error(), nil)
		return
	}

	internal.Respond(c, 200, true, "Xóa khóa học thành công", nil)
	// check if user is admin
	// if user.UserRole != models.RoleAdmin {
	// 	internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
	// 	return
	// }

	// course := models.Course{ID: courseID}

	// // check if course exists
	// if err := app.DB.First(&course).Error; err != nil {
	// 	internal.Respond(c, 404, false, "Không tìm thấy khóa học", nil)
	// 	return
	// }

	// // delete course
	// if err := app.DB.Delete(&course).Error; err != nil {
	// 	internal.Respond(c, 500, false, err.Error(), nil)
	// 	return
	// }

	// internal.Respond(c, 200, true, "Delete course success", nil)
}
