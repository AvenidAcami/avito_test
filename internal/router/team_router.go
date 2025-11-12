package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTeamRoutes(r *gin.Engine, db *gorm.DB) {
	TeamsGroup := r.Group("/team")
	{
		TeamsGroup.POST("/add")
		TeamsGroup.GET("/get")
	}
}
