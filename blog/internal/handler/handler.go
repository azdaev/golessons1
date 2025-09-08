package handler

import (
	"log"
	"net/http"
	"strconv"

	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req service.CreateUserReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println("Error in CreateUser handler: %w", err)
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	err = h.service.CreateUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Попробуйте позже")
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetUser(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		log.Println("Error in handler.GetUser, strconv.Atoi: %w", err)
		c.JSON(http.StatusBadRequest, "Невалидные данные, пользователь не найден")
		return
	}

	user, err := h.service.GetUser(c, id)
	if err != nil {
		log.Println("Error in handler.GetUser, service.GetUser: %w", err)
		c.JSON(http.StatusBadRequest, "Такой пользователь не найден")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":     user.Name,
		"is_admin": user.IsAdmin,
		"email":    user.Email,
	})
}
