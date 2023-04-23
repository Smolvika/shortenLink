package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortenLink/pkg/validite"
)

type LinkInput struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	Date        string `json:"date" binding:"required"`
}

func (h *Handler) CreateShortUrl(c *gin.Context) {
	var input LinkInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validite.OriginalURL(input.OriginalUrl); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validite.Date(input.Date); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	shortURL, err := h.services.CreateShortUrl(input.OriginalUrl, input.Date)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"shortURL": shortURL,
	})
}

func (h *Handler) GetOriginalUrl(c *gin.Context) {
	url := c.Param("short")

	if err := validite.ShortURL(url); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	originalURL, err := h.services.GetShortUrl(url)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if originalURL == "" {
		newErrorResponse(c, http.StatusNotFound, "not such ShortURL")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"originalURL": originalURL,
	})
}
