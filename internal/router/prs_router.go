package router

import (
	"avito_test/internal/repository"

	"github.com/gin-gonic/gin"
)

func InitPRsRoutes(r *gin.Engine, baseRepo repository.BaseRepository) {
	prsGroup := r.Group("/pullRequest")
	{
		prsGroup.POST("/create")
		prsGroup.POST("/merge")
		prsGroup.POST("/reassign")
	}

}
