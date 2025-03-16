package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lahaehae/crud_project/internal/service"
)

type SubUserHandler struct {
    service *service.SubUserService
}

func NewSubUserHandler(service *service.SubUserService) *SubUserHandler {
    return &SubUserHandler{
        service: service,
    }
}

func (h *SubUserHandler) CreateSubUser(c *gin.Context) {
    // Получаем ID владельца из контекста (установленного middleware)
    ownerID, _ := c.Get("user_id")

    var request struct {
        Name  string `json:"name" binding:"required"`
        Email string `json:"email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    subUser, err := h.service.CreateSubUser(c, ownerID.(string), request.Name, request.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, subUser)
}

func (h *SubUserHandler) GetUserSubUsers(c *gin.Context) {
    ownerID, _ := c.Get("user_id")

    subUsers, err := h.service.GetUserSubUsers(c, ownerID.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, subUsers)
}

// Другие обработчики...