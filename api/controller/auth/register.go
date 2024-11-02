package auth

import (
	"be/api/controller/users"
	"be/bootstrap"
	"be/db/sqlc"
	"be/internal"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username     string `json:"username" binding:"required,alphanum"`
	Password     string `json:"password" binding:"required,min=6"`
	UserFullname string `json:"user_fullname" binding:"required"`
	UserRole     string `json:"user_role" binding:"required"`
	Year         int32  `json:"year" binding:"required"`
	UserEmail    string `json:"user_email" binding:"required,email"`
}

// Register user
// @Summary Register user
// @Description Register user
// @Tags Auth
// @Accept json
// @Produce json
// @Param RegisterRequest body controller.RegisterRequest true "RegisterRequest"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /register [post]
func Register(c *gin.Context, app *bootstrap.App) {
	r := RegisterRequest{}

	if err := c.ShouldBindJSON(&r); err != nil {
		internal.Respond(c, 400, false, err.Error(), nil)
		return
	}

	// if user want to register as admin
	if sqlc.Role(r.UserRole) == sqlc.RoleAdmin || sqlc.Role(r.UserRole) == sqlc.RoleLecturer {
		sess, exists := c.Get("session")

		if !exists {
			internal.Respond(c, 401, false, "Unauthorized", nil)
			return
		}

		session := sess.(sqlc.Session)
		user, err := app.DB.GetUserById(c, session.UserID)

		if err != nil {
			app.Logger.Error().Msg(err.Error())
			internal.Respond(c, 500, false, "Internal server error", nil)
			return
		}

		if user.UserRole != sqlc.RoleAdmin {
			internal.Respond(c, 401, false, "Permission denied", nil)
			return
		}
	}

	user, err := app.DB.ValidateNewUser(c, sqlc.ValidateNewUserParams{
		Username:  r.Username,
		UserEmail: r.UserEmail,
	})

	if err == nil {
		internal.Respond(c, 403, false, "User existed", nil)
		return
	}

	hashedPassword, err := internal.HashPassword(r.Password)

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	app.Logger.Info().Msg(r.UserRole)

	user, err = app.DB.CreateUser(c, sqlc.CreateUserParams{
		UserID:       internal.GenerateUUID(),
		Username:     r.Username,
		UserFullname: r.UserFullname,
		UserEmail:    r.UserEmail,
		Password:     hashedPassword,
		UserRole:     sqlc.Role(r.UserRole),
		Year:         r.Year,
		CreatedAt:    internal.GetCurrentTime(),
		UpdatedAt:    internal.GetCurrentTime(),
	})

	if err != nil {
		app.Logger.Error().Err(err).Msg(err.Error())
		internal.Respond(c, 500, false, "Internal server error", nil)
		return
	}

	internal.Respond(c, 200, true, "User created", users.UserInfoResponse{
		UserID:       user.UserID,
		Username:     user.Username,
		UserEmail:    user.UserEmail,
		UserFullname: user.UserFullname,
		UserRole:     user.UserRole,
		Year:         user.Year,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})
}
