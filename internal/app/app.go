package app

import (
	"avito_test/config"
	"avito_test/internal/repository"
	"avito_test/internal/router"

	"github.com/gin-gonic/gin"
)

// TODO: разобраться зачем тут member, если есть user (делать что-то сонным - плохая идея)
// TODO: в GetReview добавить возврат ассигнутных для пр пользователей
func Run() {
	r := gin.Default()

	config.InitENV()
	db := config.InitDb()

	baseRepo := repository.NewBaseRepository(db)
	{
		router.InitTeamRoutes(r, baseRepo)
		router.InitUserRoutes(r, baseRepo)
	}

	r.Run(":8080")
}
