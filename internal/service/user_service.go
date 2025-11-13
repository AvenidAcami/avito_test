package service

import (
	"avito_test/internal/model"
	"avito_test/internal/repository"
)

type UserService struct {
	repo repository.IUserRepository
}

type IUserService interface {
	SetlsActive(string, bool) (model.User, error)
	GetReview(string) ([]model.PullRequest, error)
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return UserService{repo: repo}
}

func (us UserService) SetlsActive(userId string, isActive bool) (model.User, error) {
	return us.repo.SetlsActive(userId, isActive)
}

func (us UserService) GetReview(userId string) ([]model.PullRequest, error) {
	return us.repo.GetReview(userId)
}
