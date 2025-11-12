package repository

import (
	"avito_test/internal/model"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

type ITeamRepository interface {
	AddTeam(model.Team) error
	GetTeam(string) (model.Team, error)
}

func NewTeamRepository(db *gorm.DB) ITeamRepository {
	return &TeamRepository{db: db}
}

func (tr *TeamRepository) AddTeam(team model.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tx := tr.db.WithContext(ctx).Begin()

	defer cancel()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tr.insertTeamInfo(tx, team.TeamName); err != nil {
		tx.Rollback()
		return err
	}

	if err := tr.insertMembersInfo(tx, team); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (tr *TeamRepository) GetTeam(teamName string) (model.Team, error) {
	var team model.Team

	isExists, err := tr.teamExists(teamName)
	if err != nil {
		return team, err
	}
	if !isExists {
		return team, errors.New("team not exitsts")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	members, err := tr.getTeamMembers(ctx, teamName)
	if err != nil {
		return team, err
	}

	team.TeamName = teamName
	team.Members = members

	return team, nil
}

func (tr *TeamRepository) teamExists(teamName string) (bool, error) {
	var count int64
	err := tr.db.Table("teams").Where("team_name = ?", teamName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (tr *TeamRepository) getTeamMembers(ctx context.Context, teamName string) ([]model.Member, error) {
	members := make([]model.Member, 0)
	err := tr.db.WithContext(ctx).Table("members").Where("team_name = ?", teamName).Find(&members).Error
	if err != nil {
		return members, err
	}
	return members, nil
}

func (tr *TeamRepository) insertTeamInfo(tx *gorm.DB, teamName string) error {
	if err := tx.Table("teams").Create(&struct{ TeamName string }{TeamName: teamName}).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"teams_pkey\"") {
			return errors.New("team already exists")
		}
		log.Println(err.Error())
		return errors.New("something went wrong")
	}
	return nil
}

func (tr *TeamRepository) insertMembersInfo(tx *gorm.DB, team model.Team) error {
	membersDb := make([]model.MemberDb, len(team.Members))
	for ind, val := range team.Members {
		membersDb[ind] = model.MemberDb{
			UserId:   val.UserId,
			Username: val.Username,
			IsActive: val.IsActive,
			TeamName: team.TeamName,
		}
	}
	if err := tx.Table("members").Create(&membersDb).Error; err != nil {
		return errors.New("something went wrong")
	}
	return nil
}
