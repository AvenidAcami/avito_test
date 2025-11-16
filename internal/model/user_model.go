package model

type User struct {
	UserId   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	TeamName string `json:"team_name" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}
