package repository

import (
	"avito_test/internal/model"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrPullRequestAlreadyExists = errors.New("pull request already exists")
	ErrAuthorOrTeamNotFound     = errors.New("author or team not found")
)

type PRsRepository struct {
	baseRepo BaseRepository
}

type IPRsRepository interface {
	Create(pRId string, pRName string, authorId string) (model.PullRequest, error)
	Merge(pRId string) (model.PullRequestWIds, error)
	Reassign(pRId, oldUserId string) (model.PullRequestWIds, string, error)
}

func NewPRsRepository(baseRepo BaseRepository) IPRsRepository {
	return &PRsRepository{baseRepo: baseRepo}
}

func (prr *PRsRepository) Create(pRId string, pRName string, authorId string) (model.PullRequest, error) {
	var pr model.PullRequest
	pr.PullRequestId = pRId
	pr.PullRequestName = pRName
	pr.AuthorId = authorId
	pr.Status = "OPEN"
	pr.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tx := prr.baseRepo.DB.WithContext(ctx).Begin()

	defer cancel()

	if err := tx.Error; err != nil {
		return pr, err
	}

	if err := tx.Table("pull_requests").Create(&pr).Error; err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				tx.Rollback()
				return pr, ErrPullRequestAlreadyExists

			case "23503":
				tx.Rollback()
				return pr, ErrAuthorOrTeamNotFound
			}
		}

		tx.Rollback()
		return pr, err
	}

	if err := prr.assignReviewers(tx, pRId); err != nil {
		tx.Rollback()
		return pr, err
	}

	return pr, tx.Commit().Error
}

func (prr *PRsRepository) assignReviewers(tx *gorm.DB, pRId string) error {
	userIds := make([]string, 0)
	tx.Table("members").
		Order("RANDOM()").
		Limit(2).
		Where("is_active = ?", true).
		Select("user_id").
		Find(&userIds)

	for _, val := range userIds {
		if err := tx.Table("members").
			Where("user_id = ?", val).
			Update("pull_request_id", pRId).
			Error; err != nil {
			return err
		}

		if err := tx.Table("members").
			Where("user_id = ?", val).
			Update("is_active", false).
			Error; err != nil {
			return err
		}
	}

	return nil
}

func (prr *PRsRepository) Merge(pRId string) (model.PullRequestWIds, error) {
	var pr model.PullRequestWIds

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tx := prr.baseRepo.DB.WithContext(ctx).Begin()

	defer cancel()

	if err := tx.Table("pull_requests").
		Where("pull_request_id = ?", pRId).
		Update("status", "MERGED").
		Error; err != nil {
		tx.Rollback()
		return pr, err
	}

	if err := tx.Table("pull_requests").
		Where("pull_request_id = ?", pRId).
		Update("mergedAt", time.Now()).
		Error; err != nil {
		tx.Rollback()
		return pr, err
	}

	pr, err := prr.getPRInfo(tx, pRId)
	if err != nil {
		return pr, err
	}

	return pr, tx.Commit().Error
}

func (prr *PRsRepository) getPRInfo(tx *gorm.DB, pRId string) (model.PullRequestWIds, error) {
	var pr model.PullRequestWIds
	var reviewers []string

	if err := tx.Table("pull_requests").
		Where("pull_request_id = ?", pRId).
		First(&pr).Error; err != nil {
		return pr, err
	}

	if err := tx.Table("members").
		Select("user_id").
		Where("pull_request_id = ?", pRId).
		Find(&reviewers).Error; err != nil {
		return pr, err
	}

	pr.AssignedReviewers = reviewers
	return pr, nil
}

func (prr *PRsRepository) Reassign(pRId, oldUserId string) (model.PullRequestWIds, string, error) {
	var replacedBy string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tx := prr.baseRepo.DB.WithContext(ctx).Begin()

	defer cancel()

	if err := tx.Table("members").
		Where("user_id = ?", oldUserId).
		Update("is_active = ?", true).Error; err != nil {
		tx.Rollback()
		return model.PullRequestWIds{}, replacedBy, err
	}

	if err := tx.Table("members").
		Where("is_active = ?", true).
		Where("user_id != ?", oldUserId).
		Select("user_id").
		First(&replacedBy).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return model.PullRequestWIds{}, replacedBy, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pr, err := prr.getPRInfo(tx, pRId)
			if err != nil {
				return pr, "", err
			}
			return pr, "", nil
		}
	}

	if err := tx.Table("members").
		Where("user_id = ?", replacedBy).
		Update("is_active = ?", false).
		Error; err != nil {
		tx.Rollback()
		return model.PullRequestWIds{}, replacedBy, err
	}

	pr, err := prr.getPRInfo(tx, pRId)
	if err != nil {
		return pr, "", err
	}

	return pr, replacedBy, tx.Commit().Error
}
