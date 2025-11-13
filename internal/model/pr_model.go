package model

type PullRequest struct {
	PullRequestId   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorId        string `json:"author_id" binding:"required"`
	Status          string `json:"status" binding:"required"`
}
