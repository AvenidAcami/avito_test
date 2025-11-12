package repository

import "gorm.io/gorm"

type TeamRepository struct {
	db *gorm.DB
}

type ITeamRepository interface{}

func NewTeamRepository(db *gorm.DB) ITeamRepository {
	return TeamRepository{db: db}
}
