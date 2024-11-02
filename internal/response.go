package internal

import (
	"be/models"

	"github.com/gin-gonic/gin"
)

func Respond(c *gin.Context, statusCode int, success bool, message string, data interface{}) {
	c.JSON(statusCode, models.Response{
		Status:  success,
		Message: message,
		Data:    data,
	})
}
