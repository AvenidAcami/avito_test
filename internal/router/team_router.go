package router

import (
	"avito_test/internal/controller"
	"avito_test/internal/repository"
	"avito_test/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTeamRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repository.NewTeamRepository(db)
	service := service.NewTeamService(repo)
	controller := controller.NewTeamController(service)

	TeamsGroup := r.Group("/team")
	{
		TeamsGroup.POST("/add", controller.AddTeam)
		TeamsGroup.GET("/get", controller.GetTeam)
	}
}
