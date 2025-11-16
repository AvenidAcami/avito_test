package router

import (
	"avito_test/internal/controller"
	"avito_test/internal/repository"
	"avito_test/internal/service"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.Engine, baseRepo repository.BaseRepository) {
	userRepository := repository.NewUserRepository(baseRepo)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	userGroup := r.Group("/user")
	{
		userGroup.GET("/setlsActive", userController.SetlsActive)
		userGroup.POST("/getReview", userController.GetReview)
	}
}
