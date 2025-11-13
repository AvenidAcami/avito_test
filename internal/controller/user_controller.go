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
		IsActive bool   `json:"is_active" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    "NOT_FOUND",
			"message": "wrong bodys",
		})
		return
	}

	// Сделать возврат обновленного user
	user, err := uc.service.SetlsActive(body.UserId, body.IsActive)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    "NOT_FOUND",
			"message": "wrong bodys",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
