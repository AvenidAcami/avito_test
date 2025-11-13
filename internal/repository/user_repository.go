package repository

import "avito_test/internal/model"

type UserRepository struct {
	baseRepo BaseRepository
}

type IUserRepository interface {
	SetlsActive(string, bool) (model.User, error)
	GetReview(string) ([]model.PullRequest, error)
}

func NewUserRepository(baseRepo BaseRepository) IUserRepository {
	return &UserRepository{baseRepo: baseRepo}
}

func (ur *UserRepository) SetlsActive(userId string, isActive bool) (model.User, error) {
	var user model.User

	err := ur.baseRepo.DB.Table("members").Update("is_active", isActive).Where("user_id = ?", userId).Error
	if err != nil {
		return user, err
	}

	err = ur.baseRepo.DB.Table("members").Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return user, nil
	}
	return user, nil
}

func (ur *UserRepository) GetReview(userId string) ([]model.PullRequest, error) {
	var prs []model.PullRequest

	err := ur.baseRepo.DB.Table("pull_requests").Where("user_id = ?", userId).Find(&prs).Error
	if err != nil {
		return prs, err
	}

	return prs, nil
}
