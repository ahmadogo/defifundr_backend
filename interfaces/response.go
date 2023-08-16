package interfaces

import "github.com/gin-gonic/gin"

func Response(status int, data interface{}) gin.H {
	return gin.H{
		"status":  status,
		"message": "success",
		"data":    data,
	}
}


func CampaignResponse(status int, data interface{}, user interface{}) gin.H {
	return gin.H{
		"status":  status,
		"message": "success",
		"data":    data,
		"user":    user,
	}
}