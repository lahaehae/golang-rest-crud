package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lahaehae/crud_project/internal/models"
	"github.com/lahaehae/crud_project/internal/repository"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

type TransferRequest struct {
	FromID  int64 `json:"from_id" binding:"required"`
	ToID    int64 `json:"to_id" binding:"required"`
	Balance int64 `json:"balance" binding:"required"`
}

// Создание пользователя
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.repo.CreateUser(c.Request.Context(), user.Name, user.Email, user.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newUser)
}

func (h *UserHandler) TransferFunds(c *gin.Context){
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transferFunds, err := h.repo.TransferFunds(c.Request.Context(), req.FromID, req.ToID, req.Balance)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction Failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer successful",
		"user": gin.H{
			"id": transferFunds.Id,
			"balance": transferFunds.Balance,
		},
	})
}

// Получение пользователя по ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.repo.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Обновление пользователя
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.repo.UpdateUser(c.Request.Context(), id, user.Name, user.Email, user.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// Удаление пользователя
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.repo.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}


