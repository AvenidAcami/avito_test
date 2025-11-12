package controller

import (
	"avito_test/internal/model"
	"avito_test/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamController struct {
	service service.ITeamService
}

func NewTeamController(service service.ITeamService) TeamController {
	return TeamController{service: service}
}

func (tc *TeamController) AddTeam(ctx *gin.Context) {
	var team model.Team

	if err := ctx.ShouldBindJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    "TEAM_EXISTS",
			"message": "something wrong with request body",
		})
		return
	}

	if err := tc.service.AddTeam(team); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    "TEAM_EXISTS",
			"message": "team already exists",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"team_name": team.TeamName,
		"members":   team.Members,
	})
}

func (tc *TeamController) GetTeam(ctx *gin.Context) {
	teamName := ctx.Query("team_name")
	team, err := tc.service.GetTeam(teamName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "team not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"team_name": team.TeamName,
		"members":   team.Members,
	})
}
