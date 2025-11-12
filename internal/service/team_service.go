package service

import (
	"avito_test/internal/repository"
)

type TeamService struct {
	repository repository.ITeamRepository
}

type ITeamService interface{}

func NewTeamRepository(repository repository.ITeamRepository) ITeamService {
	return TeamService{repository: repository}
}
