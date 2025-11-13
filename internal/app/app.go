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
	// TODO: сделать функцию создания базе репозитори
	baseRepo := repository.BaseRepository{DB: db}
	{
		router.InitTeamRoutes(r, baseRepo)
	}

	r.Run(":8080")
}
