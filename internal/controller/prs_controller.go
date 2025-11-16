package controller

import (
	"avito_test/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrPullRequestAlreadyExists = errors.New("pull request already exists")

type PRsController struct {
	service service.IPRsService
}

func NewPRsController(service service.IPRsService) PRsController {
	return PRsController{service: service}
}

func (prc *PRsController) Create(ctx *gin.Context) {
	var body struct {
		PullRequestId   string `json:"pull_request_id" binding:"required"`
		PullRequestName string `json:"pull_request_name" binding:"required"`
		AuthorId        string `json:"author_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "author or team not found",
		})
		return
	}

	pr, err := prc.service.Create(body.PullRequestId, body.PullRequestName, body.AuthorId)
	if err != nil {
		if errors.Is(err, ErrPullRequestAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{
				"code":    "PR_EXISTS",
				"message": "pull request already exists",
			})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "author or team not found",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"pr": pr,
	})
}

func (prc *PRsController) Merge(ctx *gin.Context) {
	var body struct {
		PullRequestId string `json:"pull_request_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "wrong body",
		})
		return
	}

	pr, err := prc.service.Merge(body.PullRequestId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "pull request not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"pr": pr,
	})
}

func (prc *PRsController) Reassign(ctx *gin.Context) {
	var body struct {
		PullRequestId string `json:"pull_request_id" binding:"required"`
		OldUserId     string `json:"old_user_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "wrong body",
		})
		return
	}

	pr, replacedBy, err := prc.service.Reassign(body.PullRequestId, body.OldUserId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "pull request or user not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": replacedBy,
	})
}
