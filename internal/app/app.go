package app

import (
	"avito_test/config"
	"avito_test/internal/router"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	config.InitENV()
	db := config.InitDb()

	{
		router.InitTeamRoutes(r, db)
	}

	r.Run(":8080")
}
