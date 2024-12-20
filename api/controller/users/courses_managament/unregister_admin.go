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

type UserUnRegisterAdminCourseRequest struct {
	UserID 	   uuid.UUID `json:"user_id" binding:"required"`
	CourseID       uuid.UUID `json:"course_id" binding:"required"`
	CourseYear     int32     `json:"course_year" binding:"required"`
	CourseSemester int32     `json:"course_semester" binding:"required"`
}

// UserUnRegisterAdminCourse docs
// @Summary User unregister admin course
// @Description User unregister admin course
// @Tags User
// @Accept json
// @Produce json
// @Param UserUnRegisterAdminCourseRequest body UserUnRegisterAdminCourseRequest true "UserUnRegisterCourseRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /user/course/unregister_admin [post]
func UserUnRegisterAdminCourse(c *gin.Context, app *bootstrap.App) {
	// if app.State != bootstrap.ACTIVE {
	// 	internal.Respond(c, 403, false, fmt.Sprintf("Máy chủ không ở trạng thái ACTIVE, trạng thái hiện tại là %s", app.State), nil)
	// 	return
	// }
	// validate request
	req := UserUnRegisterAdminCourseRequest{}
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

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
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

	if err := app.DB.Where("id = ?", req.UserID).First(&models.User{}).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
	}

	// check if user has already registered the course
	// registeredCourse := models.RegisteredCourse{UserID: user.ID, CourseID: course.ID}

	// register course
	registeredCourse := models.RegisteredCourse{
		UserID:         req.UserID,
		CourseID:       course.ID,
		CourseYear:     req.CourseYear,
		CourseSemester: req.CourseSemester,
		Status:         models.CoStatusProgressing,
	}

	if err := app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(&registeredCourse).First(&registeredCourse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				internal.Respond(c, 404, false, "Không tìm thấy khóa học", nil)
				return err
			}

			app.Logger.Error().Err(err).Msg(err.Error())
			return err
		}

		if err := tx.Where(&registeredCourse).Delete(&registeredCourse).Error; err != nil {
			app.Logger.Error().Err(err).Msg(err.Error())
			return err
		}

		course.CurrentEnroller = course.CurrentEnroller - 1
		// update course current enroller
		if err := tx.Save(&course).Error; err != nil {
			app.Logger.Error().Err(err).Msg(err.Error())
			return err
		}

		return nil
	}); err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Hủy đăng ký khóa học thành công", nil)
}
