package router

import (
	"avito_test/internal/controller"
	"avito_test/internal/repository"
	"avito_test/internal/service"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.Engine, baseRepo repository.BaseRepository) {
	repo := repository.NewUserRepository(baseRepo)
	service := service.NewUserService(repo)
	controller := controller.NewUserController(service)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/setIsActive", controller.SetlsActive)
		userGroup.GET("/getReview", controller.GetReview)
	}
}
