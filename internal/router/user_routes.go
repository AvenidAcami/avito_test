package router

import "github.com/gin-gonic/gin"

func InitUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.GET("/setlsActive")
		userGroup.POST("/getReview")
	}
}
