package repository

import (
	"avito_test/internal/model"
	"context"
	"time"
)

type UserRepository struct {
	baseRepo BaseRepository
}

type IUserRepository interface {
	SetlsActive(string, bool) (model.User, error)
	GetReview(string) ([]model.PullRequestWIds, error)
}

func NewUserRepository(baseRepo BaseRepository) IUserRepository {
	return &UserRepository{baseRepo: baseRepo}
}

func (ur *UserRepository) SetlsActive(userId string, isActive bool) (model.User, error) {
	var user model.User

	err := ur.baseRepo.DB.Table("members").Where("user_id = ?", userId).Update("is_active", isActive).Error
	if err != nil {
		return user, err
	}

	err = ur.baseRepo.DB.Table("members").Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) GetReview(userId string) ([]model.PullRequestWIds, error) {
	var prs []model.PullRequestWIds
	var assignedUsers []string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tx := ur.baseRepo.DB.WithContext(ctx).Begin()

	defer cancel()

	err := tx.Table("pull_requests").Where("user_id = ?", userId).Find(&prs).Error
	if err != nil {
		tx.Rollback()
		return prs, err
	}

	for ind, val := range prs {
		if err := tx.Table("members").Where("pull_request_id = ?", val.PullRequestId).Select("user_id").Find(&assignedUsers).Error; err != nil {
			tx.Rollback()
			return prs, err
		}
		prs[ind].AssignedReviewers = assignedUsers
	}

	return prs, tx.Commit().Error
}
