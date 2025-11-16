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

	if err := prr.assignReviewers(tx, pRId, authorId); err != nil {
		tx.Rollback()
		return pr, err
	}

	return pr, tx.Commit().Error
}

func (prr *PRsRepository) assignReviewers(tx *gorm.DB, pRId, authorId string) error {
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
			Where("team_name = ?", authorId).
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
	var prwid model.PullRequestWIds

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tx := prr.baseRepo.DB.WithContext(ctx).Begin()

	defer cancel()

	if err := tx.Table("pull_requests").
		Where("pull_request_id = ?", pRId).
		Update("status", "MERGED").
		Error; err != nil {
		tx.Rollback()
		return prwid, err
	}

	if err := tx.Table("pull_requests").
		Where("pull_request_id = ?", pRId).
		Update("mergedAt", time.Now()).
		Error; err != nil {
		tx.Rollback()
		return prwid, err
	}

	prwid, err := prr.getPRInfo(tx, pRId)
	if err != nil {
		return prwid, err
	}

	for _, val := range prwid.AssignedReviewers {
		if err := tx.Table("members").
			Where("user_id = ?", val).
			Update("is_active", true).
			Error; err != nil {
			return prwid, err
		}

		if err := tx.Table("members").
			Where("user_id = ?", val).
			Update("pull_request_id", nil).
			Error; err != nil {
			return prwid, err
		}
	}

	return prwid, tx.Commit().Error
}

func (prr *PRsRepository) getPRInfo(tx *gorm.DB, pRId string) (model.PullRequestWIds, error) {
	var pr model.PullRequest
	var prwid model.PullRequestWIds
	var reviewers []string

	if err := tx.Table("pull_requests").
		Where("pull_request_id = ?", pRId).
		Select("pull_request_id", "pull_request_name", "author_id", "status", "createdAt", "mergedAt").
		First(&pr).Error; err != nil {
		tx.Rollback()
		return prwid, err
	}

	if err := tx.Table("members").
		Select("user_id").
		Where("pull_request_id = ?", pRId).
		Find(&reviewers).Error; err != nil {
		tx.Rollback()
		return prwid, err
	}

	prwid.PullRequestId = pr.PullRequestId
	prwid.PullRequestName = pr.PullRequestName
	prwid.AuthorId = pr.AuthorId
	prwid.Status = pr.Status
	prwid.CreatedAt = pr.CreatedAt
	prwid.MergedAt = pr.MergedAt
	prwid.AssignedReviewers = reviewers
	return prwid, nil
}

func (prr *PRsRepository) Reassign(pRId, oldUserId string) (model.PullRequestWIds, string, error) {
	var replacedBy string
	var IsActive bool

	if err := prr.baseRepo.DB.
		Table("members").
		Where("user_id = ?", oldUserId).
		Select("is_active").
		Pluck("is_active", &IsActive).
		Error; err != nil {
		return model.PullRequestWIds{}, replacedBy, err
	}

	if IsActive {
		return model.PullRequestWIds{}, replacedBy, errors.New("user is not assigned to any pull request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tx := prr.baseRepo.DB.WithContext(ctx).Begin()

	defer cancel()

	if err := tx.Table("members").
		Where("user_id = ?", oldUserId).
		Update("is_active", true).Error; err != nil {
		tx.Rollback()
		return model.PullRequestWIds{}, replacedBy, err
	}

	if err := tx.Table("members").
		Where("user_id = ?", oldUserId).
		Update("pull_request_id", nil).
		Error; err != nil {
		return model.PullRequestWIds{}, replacedBy, err
	}

	if err := tx.Table("members").
		Where("is_active = ?", true).
		Where("user_id != ?", oldUserId).
		Select("user_id").
		Pluck("user_id", &replacedBy).Error; err != nil {
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
		Update("is_active", false).
		Error; err != nil {
		tx.Rollback()
		return model.PullRequestWIds{}, replacedBy, err
	}

	if err := tx.Table("members").
		Where("user_id = ?", replacedBy).
		Update("pull_request_id", pRId).
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
