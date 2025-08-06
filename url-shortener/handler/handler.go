package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"url-shortener-1/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const HostURL = "127.0.0.1:8085/"

type Handler struct {
	LinksRepository *repo.Repository
}

type CreateLinkRequest struct {
	Link string `json:"link"`
}

type LinkResponse struct {
	LongLink  string `json:"long_link"`
	ShortLink string `json:"short_link"`
}

func NewHandler(linksRepo *repo.Repository) Handler {
	return Handler{
		LinksRepository: linksRepo,
	}
}

func (h *Handler) CreateLink(c *gin.Context) {
	var req CreateLinkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	// Проверка что такую длинную ссылку уже сокращали
	// Если уже есть короткая ссылка -> возвращаем ее клиенту
	// Если нет -> генерируем

	existingShortLink, err := h.LinksRepository.GetShortByLong(c, req.Link)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"short": HostURL + existingShortLink,
			"long":  req.Link,
		})
		return
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println("Ошибка БД: ", err)
		c.JSON(http.StatusInternalServerError, "Ошибка базы данных")
		return
	}

	// генерация уникальной короткой ссылки

	shortLink := ""
	for {
		b := make([]byte, 6)
		rand.Read(b)
		shortLink = base64.URLEncoding.EncodeToString(b)[:6]

		isExists, err := h.LinksRepository.IsShortExists(c, shortLink)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Произошла ошибка БД, попробуйте позже")
			return
		}

		if !isExists {
			break
		}
	}

	err = h.LinksRepository.CreateLink(c, req.Link, shortLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short": HostURL + shortLink,
		"long":  req.Link,
	})
}

func (h *Handler) Redirect(c *gin.Context) {
	shortLink := c.Param("path")
	var longLink string

	row := h.db.QueryRow(c, "select long_link from links where short_link=$1", shortLink)
	err := row.Scan(&longLink)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, "Ссылка не найдена")
			return
		}
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, longLink)
}
