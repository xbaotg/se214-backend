package global

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateEnvFile(key, value string) error {
    // Read the .env file
    input, err := ioutil.ReadFile(".env")
    if err != nil {
        return err
    }

    lines := strings.Split(string(input), "\n")
    
    // Find and replace the line containing the key
    for i, line := range lines {
        if strings.HasPrefix(line, key+"=") {
            lines[i] = key + "='" + value + "'"
            break
        }
    }
    
    // Join the lines back together
    output := strings.Join(lines, "\n")
    
    // Write the modified content back to .env
    err = ioutil.WriteFile(".env", []byte(output), 0644)
    if err != nil {
        return err
    }
    
    return nil
}


// Set global state
// @Summary Set global state
// @Description Set global state
// @Tags Global
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param state query string true "State" Enums(freeze, active, setup, done)
// @Router /global/state [post]
func SetGlobalState(c *gin.Context, app *bootstrap.App) {
	sess, _ := c.Get("session")
	state := c.Query("state")
	session := sess.(models.Session)

	// get user info
	user := models.User{
		ID: session.UserID,
	}
	if err := app.DB.First(&user).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	if user.UserRole != models.RoleAdmin {
		internal.Respond(c, 403, false, "Không có quyền truy cập", nil)
		return
	}


	app.State = bootstrap.State(state)
	UpdateEnvFile("APP_INITIAL_STATE", string(app.State))
	internal.Respond(c, 200, true, "Cập nhật trạng thái thành công", nil)
}
