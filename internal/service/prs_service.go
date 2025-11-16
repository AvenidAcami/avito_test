package service

import (
	"avito_test/internal/model"
	"avito_test/internal/repository"
)

type PRsService struct {
	repo repository.IPRsRepository
}

type IPRsService interface {
	Create(pRId string, pRName string, authorId string) (model.PullRequest, error)
	Merge(pRId string) (model.PullRequestWIds, error)
	Reassign(pRId, oldUserId string) (model.PullRequestWIds, string, error)
}

func NewPRsService(repo repository.IPRsRepository) IPRsService {
	return &PRsService{repo: repo}
}

func (prs *PRsService) Create(pRId string, pRName string, authorId string) (model.PullRequest, error) {
	return prs.repo.Create(pRId, pRName, authorId)
}

func (prs *PRsService) Merge(pRId string) (model.PullRequestWIds, error) {
	return prs.repo.Merge(pRId)
}
func (prs *PRsService) Reassign(pRId, oldUserId string) (model.PullRequestWIds, string, error) {
	return prs.repo.Reassign(pRId, oldUserId)
}
