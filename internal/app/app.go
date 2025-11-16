package app

import (
	"avito_test/config"
	"avito_test/internal/repository"
	"avito_test/internal/router"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	config.InitENV()
	db := config.InitDb()

	baseRepo := repository.NewBaseRepository(db)
	{
		router.InitTeamRoutes(r, baseRepo)
		router.InitUserRoutes(r, baseRepo)
		router.InitPRsRoutes(r, baseRepo)
	}

	r.Run(":8080")
}
