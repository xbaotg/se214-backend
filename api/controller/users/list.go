package users

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListUsersRequest struct {
	// Role string `json:"role"` // admin, user, lecturer, all
}

type ListUsersResponse struct {
	UserID       string      `json:"id"`
	Username     string      `json:"username"`
	UserEmail    string      `json:"email"`
	UserFullname string      `json:"user_fullname"`
	UserRole     models.Role `json:"user_role"`
	Year         int32       `json:"year"`
}

// List users
// @Summary List users
// @Description List users
// @Tags User
// @Produce json
// @Param role query string false "Role" Enums(admin, user, lecturer, all)
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/list [get]
func ListUsers(c *gin.Context, app *bootstrap.App) {
	// req := ListUsersRequest{}
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	internal.Respond(c, 400, false, err.Error(), nil)
	// 	return
	// }

	role := c.Query("role")
	sess, _ := c.Get("session")
	session := sess.(models.Session)
	user := models.User{ID: session.UserID}

	if err := app.DB.First(&user).Error; err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	if user.UserRole == models.RoleLecturer {
		internal.Respond(c, 200, true, "List users successfully", []ListUsersResponse{
			ListUsersResponse{
				UserID:       user.ID.String(),
				Username:     user.Username,
				UserEmail:    user.UserEmail,
				UserFullname: user.UserFullname,
				UserRole:     user.UserRole,
				Year:         user.Year,
			},
		})
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Forbidden", nil)
		return
	}

	users := []models.User{}
	if role == "" {
		if err := app.DB.Find(&users).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				internal.Respond(c, 404, false, "No users found", nil)
				return
			}

			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}
	} else {
		if err := app.DB.Where("user_role = ?", models.Role(role)).Find(&users).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				internal.Respond(c, 404, false, "No users found", nil)
				return
			}

			app.Logger.Error().Err(err).Msg(err.Error())
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}
	}

	res := []ListUsersResponse{}
	for _, u := range users {
		res = append(res, ListUsersResponse{
			UserID:       u.ID.String(),
			Username:     u.Username,
			UserEmail:    u.UserEmail,
			UserFullname: u.UserFullname,
			UserRole:     u.UserRole,
			Year:         u.Year,
		})
	}

	internal.Respond(c, 200, true, "List users successfully", res)
}
