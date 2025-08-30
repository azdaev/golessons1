package handler

import (
	"errors"
	"log"
	"net/http"
	"url-shortener-1/manager"
	"url-shortener-1/model"
	"url-shortener-1/service"

	"github.com/gin-gonic/gin"
)

const (
	HostURL = "127.0.0.1:8080/"
)

type Handler struct {
	linksManager manager.LinksManager
	linksService service.LinksService
}

func New(linksManager manager.LinksManager, linksService service.LinksService) Handler {
	return Handler{
		linksManager: linksManager,
		linksService: linksService,
	}
}

func (h *Handler) CreateLink(c *gin.Context) {
	var req model.CreateLinkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		return
	}

	shortLink, err := h.linksService.CreateShortLink(c, req.Link, req.CustomShortLink)
	if err != nil {
		if errors.Is(err, service.ErrorLinkAlreadyExists) ||
			errors.Is(err, service.ErrorLinkTooShort) ||
			errors.Is(err, service.ErrorInvalidSymbolInLink) {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		log.Printf("error linksService.CreateShortLink: %v", err)
		c.JSON(http.StatusInternalServerError, "Ошибка! Попробуйте позже")
	}

	c.JSON(http.StatusOK, model.LinkResponse{
		ShortLink: HostURL + shortLink,
		LongLink:  req.Link,
	})
}

func (h *Handler) Redirect(c *gin.Context) {
	shortLink := c.Param("path")

	longLink, err := h.linksManager.GetLongByShort(c, shortLink)
	if err != nil {
		if errors.Is(err, errors.New("error link not found")) {
			c.JSON(http.StatusNotFound, "Ссылка не найдена")
			return
		}
		log.Println("error linksManager.GetLongByShort: ", err)
		c.JSON(http.StatusInternalServerError, "Произошла ошибка, попробуйте позже")
		return
	}

	err = h.linksManager.StoreRedirect(c, model.StoreRedirectParams{
		UserAgent: c.GetHeader("User-Agent"),
		LongLink:  longLink,
		ShortLink: shortLink,
	})
	if err != nil {
		log.Printf("Ошибка при StoreRedirect: %v", err)
	}

	c.Redirect(http.StatusTemporaryRedirect, longLink)
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	shortLink := c.Param("path")

	redirects, err := h.linksManager.GetRedirectsByShortLink(c, shortLink)
	if err != nil {
		log.Printf("error GetRedirectsByShortLink: %v\n", err)
		c.JSON(http.StatusInternalServerError, "Не удалось получить аналитику")
		return
	}

	c.JSON(http.StatusOK, model.AnalyticsResponse{
		TotalRedirects: len(redirects),
		Redirects:      redirects,
	})
}
