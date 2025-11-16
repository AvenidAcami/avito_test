package model

import "time"

type PullRequest struct {
	PullRequestId   string    `json:"pull_request_id" binding:"required"`
	PullRequestName string    `json:"pull_request_name" binding:"required"`
	AuthorId        string    `json:"author_id" binding:"required"`
	Status          string    `json:"status" binding:"required"`
	CreatedAt       time.Time `json:"createdAt"`
	MergedAt        time.Time `json:"mergedAt"`
}

type PullRequestWIds struct {
	PullRequestId     string    `json:"pull_request_id" binding:"required"`
	PullRequestName   string    `json:"pull_request_name" binding:"required"`
	AuthorId          string    `json:"author_id" binding:"required"`
	Status            string    `json:"status" binding:"required"`
	AssignedReviewers []string  `json:"assigned_reviewers" binding:"required"`
	CreatedAt         time.Time `json:"createdAt"`
	MergedAt          time.Time `json:"mergedAt"`
}
