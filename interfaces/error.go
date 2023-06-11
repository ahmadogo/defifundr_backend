package interfaces

import "github.com/gin-gonic/gin"

func ErrorResponse(err error, status int) gin.H {
	return gin.H{
		"status":  status,
		"message": err.Error(),
	}
}
