package internal

import (
	"be/model"

	"github.com/gin-gonic/gin"
)

func Respond(c *gin.Context, statusCode int, success bool, message string, data interface{}) {
	c.JSON(statusCode, model.Response{
		Status:  success,
		Message: message,
		Data:    data,
	})
}
