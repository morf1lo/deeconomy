package v1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/morf1lo/deeconomy-bot-api/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	services *service.Service
}

func NewHandler(logger *zap.Logger, services *service.Service) *Handler {
	return &Handler{
		logger: logger,
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()

	r.Use(cors.Default())

	v1 := r.Group("/api/v1")
	{
		oauth2 := v1.Group("/oauth2")
		{
			oauth2.GET("/signin", h.oauth2SignIn)
			oauth2.GET("/callback", h.oauth2Authorize)
		}
	}

	return r
}
