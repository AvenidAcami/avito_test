package controller

import (
	"avito_test/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.IUserService
}

func NewUserController(service service.IUserService) UserController {
	return UserController{service: service}
}

func (uc *UserController) SetlsActive(ctx *gin.Context) {
	var body struct {
		UserId   string `json:"user_id" binding:"required"`
		IsActive *bool  `json:"is_active" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    "NOT_FOUND",
			"message": "wrong body",
		})
		return
	}

	user, err := uc.service.SetlsActive(body.UserId, *body.IsActive)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "user not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (uc *UserController) GetReview(ctx *gin.Context) {
	userId := ctx.Query("user_id")

	prs, err := uc.service.GetReview(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "user not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":       userId,
		"pull_requests": prs,
	})
}
