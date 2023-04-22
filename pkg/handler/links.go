package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shortenLink/pkg/valid"
)

type LinkInput struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	ShortUrl    string
	Date        string `json:"date" binding:"required"`
}

func (h *Handler) CreateShortUrl(c *gin.Context) {
	var input LinkInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := valid.ValidateOriginalURL(input.OriginalUrl); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//if err := valid.ValidateDate(input.Date); err != nil {
	//	newErrorResponse(c, http.StatusBadRequest, err.Error())
	//	return
	//}
	shortURL, err := h.services.CreateShortUrl(input.OriginalUrl, input.ShortUrl, input.Date)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"shortURL": shortURL,
	})
}

type shortLinkInput struct {
	Url string `json:"short_url" binding:"required"`
}

func (h *Handler) GetOriginalUrl(c *gin.Context) {
	url := c.Param("short")

	if err := valid.ValidateShortURL(url); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	originalURL, err := h.services.GetShortUrl(url)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if originalURL == "" {
		newErrorResponse(c, http.StatusInternalServerError, "not such token")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"originalURL": originalURL,
	})
}
