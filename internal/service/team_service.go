package service

import (
	"avito_test/internal/model"
	"avito_test/internal/repository"
)

type TeamService struct {
	repository repository.ITeamRepository
}

type ITeamService interface {
	AddTeam(model.Team) error
	GetTeam(string) (model.Team, error)
}

func NewTeamService(repository repository.ITeamRepository) ITeamService {
	return &TeamService{repository: repository}
}

func (ts *TeamService) AddTeam(team model.Team) error {
	return ts.repository.AddTeam(team)
}

func (ts *TeamService) GetTeam(teamName string) (model.Team, error) {
	return ts.repository.GetTeam(teamName)
}
