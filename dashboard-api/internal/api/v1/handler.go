package v1

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/morf1lo/deeconomy-bot-api/internal/lib"
	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"github.com/morf1lo/deeconomy-bot-api/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			oauth2.GET("/refresh", h.oauth2Refresh)
		}

		users := v1.Group("/users")
		{
			users.GET("/@me", h.AuthMiddleware, h.usersMe)
			users.GET("/guilds", h.AuthMiddleware, h.usersGuilds)
		}
	}

	return r
}

func (h *Handler) getUserDataFromTokenClaims(ctx context.Context, accessToken string) (*model.User, error) {
	decodedToken, err := lib.DecodeAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(decodedToken["id"].(string))
	if err != nil {
		return nil, errIDIsNotValid
	}

	user, err := h.services.User.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) getUser(c *gin.Context) *model.User {
	userReq, _ := c.Get("user")

	user, ok := userReq.(model.User)
	if !ok {
		return nil
	}

	return &user
}
