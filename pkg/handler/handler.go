package handler

import (
	"github.com/gin-gonic/gin"
	"shortenLink/pkg/service"
)

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	links := router.Group("/url")
	{
		links.POST("/short", h.CreateShortUrl)
		links.GET("/:short", h.GetOriginalUrl)
	}

	return router
}
