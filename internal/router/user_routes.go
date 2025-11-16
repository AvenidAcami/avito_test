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

	userGroup := r.Group("/user")
	{
		userGroup.GET("/setlsActive", controller.SetlsActive)
		userGroup.POST("/getReview", controller.GetReview)
	}
}
