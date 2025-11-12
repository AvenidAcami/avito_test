package model

type Team struct {
	TeamName string   `json:"team_name" binding:"required"`
	Members  []Member `json:"members" binding:"required"`
}

type Member struct {
	UserId   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type MemberDb struct {
	UserId   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
	TeamName string `json:"team_name" binding:"required"`
}
