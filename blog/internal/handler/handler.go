package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"blog/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers(c)
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

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.service.GetPosts(c)
	if err != nil {
		log.Println("get posts handler: %w", err)
		c.JSON(http.StatusInternalServerError, "Не удалось получить посты")
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPostsByUserID(c *gin.Context) {
	param := c.Param("id")
	userId, err := strconv.Atoi(param)
	if err != nil {
		log.Println("Error in hanlder GetPostsByUserId problem with id: %w", err)
		c.JSON(http.StatusBadRequest, "Некорректный айди")
	}

	posts, err := h.service.GetPostsByUserID(c, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Error in GetPostsByUserId with no rows: %w", err)
			c.JSON(http.StatusBadRequest, "Посты не найдены")
			return
		}
		log.Println("Error in GetPostsByUserId: %w", err)
		c.JSON(http.StatusBadGateway, "Произошла какая-то ошибка, повторите позже")
		return
	}

	c.JSON(http.StatusOK, posts)
}
