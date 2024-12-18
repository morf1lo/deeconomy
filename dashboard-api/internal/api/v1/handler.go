package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morf1lo/deeconomy-bot-api/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	services *service.Service
	httpClient *http.Client
}

func New(logger *zap.Logger, services *service.Service) *Handler {
	return &Handler{
		logger: logger,
		services: services,
		httpClient: &http.Client{},
	}
}

func InitRoutes() {
	r := gin.New()

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/@me")
		}

		guilds := v1.Group("/guilds")
		{
			guilds.GET("")
		}
	}
}
