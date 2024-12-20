package subject

import (
	"be/bootstrap"
	"be/internal"
	"be/models"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
	// "github.com/google/uuid"
)

type AddPrerequisiteRequest struct {
	CourseID       string `json:"course_id" binding:"required"`
	PrerequisiteID string `json:"prerequisite_id" binding:"required"`
}

// Add Prerequisite docs
// @Summary Add Prerequisite
// @Description Add Prerequisite
// @Tags Subject
// @Accept json
// @Produce json
// @Param AddPrerequisiteRequest body AddPrerequisiteRequest true "AddPrerequisiteRequest"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /subject/add_prerequisite [post]
func AddPrerequisite(c *gin.Context, app *bootstrap.App) {
	req := AddPrerequisiteRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}

	if req.CourseID == req.PrerequisiteID {
		internal.Respond(c, 400, false, "Môn học không thể tự thêm môn tiên quyết cho chính nó", nil)
		return
	}

	course := models.AllCourses{CourseName: req.CourseID}
	if err := app.DB.First(&course).Error; err != nil {
		internal.Respond(c, 404, false, "Môn học không tồn tại", nil)
		return
	}

	if !course.Status {
		internal.Respond(c, 400, false, "Môn học không hoạt động", nil)
		return
	} 

	prerequisite := models.AllCourses{CourseName: req.PrerequisiteID}
	if err := app.DB.First(&prerequisite).Error; err != nil {
		internal.Respond(c, 404, false, "Không tìm thấy môn học tiên quyết", nil)
		return
	}

	if !prerequisite.Status {
		internal.Respond(c, 400, false, "Không thể thêm môn học tiên quyết không hoạt động", nil)
		return
	}

	prerequisiteCourse := models.PrerequisiteCourse{
		CourseID:       req.CourseID,
		PrerequisiteID: req.PrerequisiteID,
	}

	if res := app.DB.Create(&prerequisiteCourse); res.Error != nil {
		app.Logger.Info().Msgf("%d", res.RowsAffected)
		if res.RowsAffected == 0 {
			internal.Respond(c, 400, false, "Môn học tiên quyết đã tồn tại", nil)
			return
		}
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}
	
	internal.Respond(c, 200, true, "Thêm môn học tiên quyết thành công", prerequisiteCourse)
}
