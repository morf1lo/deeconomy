package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/morf1lo/deeconomy-bot-api/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	services *service.Service
}

func New(logger *zap.Logger, services *service.Service) *Handler {
	return &Handler{
		logger: logger,
		services: services,
	}
}

func (h *Handler) InitRoutes() {
	r := gin.New()

	v1 := r.Group("/api/v1")
	{
		oauth2 := v1.Group("/oauth2")
		{
			oauth2.GET("/signin", h.oauth2SignIn)
			oauth2.GET("/callback", h.oauth2Authorize)
		}
	}
}
