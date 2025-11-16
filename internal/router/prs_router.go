package router

import (
	"avito_test/internal/controller"
	"avito_test/internal/repository"
	"avito_test/internal/service"

	"github.com/gin-gonic/gin"
)

func InitPRsRoutes(r *gin.Engine, baseRepo repository.BaseRepository) {
	repo := repository.NewPRsRepository(baseRepo)
	service := service.NewPRsService(repo)
	controller := controller.NewPRsController(service)

	prsGroup := r.Group("/pullRequest")
	{
		prsGroup.POST("/create", controller.Create)
		prsGroup.POST("/merge", controller.Merge)
		prsGroup.POST("/reassign", controller.Reassign)
	}

}
