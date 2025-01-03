package global

import (
	"be/bootstrap"
	"be/internal"
	"be/models"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

func readSQLFile(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read all contents of the file
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Return the content as a string
	return string(content), nil
}

// UploadData godoc
// @Summary Upload data
// @Description Upload data
// @Tags Global
// @Produce json
// @Success 200 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param file formData file true "File"
// @Router /global/upload_data [post]
func UploadData(c *gin.Context, app *bootstrap.App) {
	if app.State != bootstrap.SETUP {
		internal.Respond(c, 403, false, "Chỉ có thể upload dữ liệu khi ứng dụng ở trạng thái SETUP", nil)
		return
	}

	sess, _ := c.Get("session")
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
	file, err := c.FormFile("file")
	if err != nil {
		internal.Respond(c, 400, false, "Hãy chọn file", nil)
		return
	}

	c.SaveUploadedFile(file, "./data/"+file.Filename)
	defer os.Remove("./data/" + file.Filename)

	// read file content
	content, err := readSQLFile("./data/" + file.Filename)
	if err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	// insert data to database
	if err := app.DB.Exec(content).Error; err != nil {
		internal.Respond(c, 500, false, "Lỗi máy chủ", nil)
		return
	}

	internal.Respond(c, 200, true, "Upload dữ liệu thành công", nil)
}
