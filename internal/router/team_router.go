package router

import (
	"avito_test/internal/controller"
	"avito_test/internal/repository"
	"avito_test/internal/service"

	"github.com/gin-gonic/gin"
)

func InitTeamRoutes(r *gin.Engine, baseRepo repository.BaseRepository) {
	repo := repository.NewTeamRepository(baseRepo)
	service := service.NewTeamService(repo)
	controller := controller.NewTeamController(service)

	TeamsGroup := r.Group("/team")
	{
		TeamsGroup.POST("/add", controller.AddTeam)
		TeamsGroup.GET("/get", controller.GetTeam)
	}
}
