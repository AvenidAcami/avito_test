package controller

import (
	"avito_test/internal/service"

	"github.com/gin-gonic/gin"
)

type TeamController struct {
	service service.ITeamService
}

func NewTeamController(service service.ITeamService) TeamController {
	return TeamController{service: service}
}

func (tc *TeamController) AddTeam(ctx *gin.Context) {

}
