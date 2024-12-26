package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/morf1lo/deeconomy-bot-api/internal/dto"
	"github.com/morf1lo/deeconomy-bot-api/internal/lib"
	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type userService struct {
	logger *zap.Logger
	repo *repository.Repository
	httpClient *http.Client
}

func newUserService(logger *zap.Logger, repo *repository.Repository) User {
	return &userService{
		logger: logger,
		repo: repo,
		httpClient: &http.Client{},
	}
}

func (s *userService) Authorize(ctx context.Context, discordCode string) (string, string, error) {
	// /oauth2/token ENDPOINT REQUEST
	oauth2Endpoint := fmt.Sprintf("%s/oauth2/token", DISCORD_HOST)

	reqBody := &url.Values{}
	reqBody.Add("client_id", os.Getenv("CLIENT_ID"))
	reqBody.Add("client_secret", os.Getenv("CLIENT_SECRET"))
	reqBody.Add("grant_type", "authorization_code")
	reqBody.Add("redirect_uri", os.Getenv("DISCORD_REDIRECT_URI"))
	reqBody.Add("code", discordCode)

	req, err := http.NewRequest(http.MethodPost, oauth2Endpoint, strings.NewReader(reqBody.Encode()))
	if err != nil {
		s.logger.Sugar().Errorf("failed to CREATE request for '/oauth2/token' discord endpoint: %s", err.Error())
		return "", "", errInternal
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	oauth2Resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Sugar().Errorf("failed to DO request to authorize user: %s", err.Error())
		return "", "", err
	}
	defer oauth2Resp.Body.Close()

	var oauth2Response dto.DiscordOAuth2Resp
	if err := json.NewDecoder(oauth2Resp.Body).Decode(&oauth2Response); err != nil {
		s.logger.Sugar().Errorf("failed to decode discord oauth2 response body: %s", err.Error())
		return "", "", errInternal
	}

	if oauth2Resp.StatusCode != http.StatusOK {
		s.logger.Sugar().Errorf("error from discord api(status: %d): %s", oauth2Resp.StatusCode, oauth2Response.Message)
		return "", "", fmt.Errorf("discord error: %s", oauth2Response.Message)
	}

	// /users/@me ENDPOINT REQUEST
	usersEndpoint := fmt.Sprintf("%s/users/@me", DISCORD_HOST)
	req, err = http.NewRequest(http.MethodGet, usersEndpoint, nil)
	if err != nil {
		s.logger.Sugar().Errorf("failed to CREATE request for '/users/@me' discord endpoint: %s", err.Error())
		return "", "", errInternal
	}

	req.Header.Add("Authorization", "Bearer " + oauth2Response.AccessToken)

	userResp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Sugar().Errorf("failed to DO request to get user info: %s", err.Error())
		return "", "", err
	}
	defer userResp.Body.Close()

	userRespBodyBytes, err := io.ReadAll(userResp.Body)
	if err != nil {
		s.logger.Sugar().Errorf("failed to read discord user response body: %s", err.Error())
		return "", "", errInternal
	}

	if userResp.StatusCode != http.StatusOK {
		var discordError dto.DiscordError
		if err := json.Unmarshal(userRespBodyBytes, &discordError); err != nil {
			s.logger.Sugar().Errorf("failed to decode discord error response: %s", err.Error())
			return "", "", fmt.Errorf("discord error: unknown error")
		}

		s.logger.Sugar().Errorf("error from discord api(status: %d, code: %d): %s", userResp.StatusCode, discordError.Code, discordError.Message)
		return "", "", fmt.Errorf("discord error: %s", discordError.Message)
	}

	var userResponse discordgo.User
	if err := json.Unmarshal(userRespBodyBytes, &userResponse); err != nil {
		s.logger.Sugar().Errorf("failed to decode discord user response: %s", err.Error())
		return "", "", errInternal
	}

	// Creating or updating user
	user, err := s.repo.Mongo.User.FindByDiscordID(ctx, userResponse.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		s.logger.Sugar().Errorf("failed to find user(%s): %s", userResponse.ID)
		return "", "", errInternal
	}
	if err == mongo.ErrNoDocuments {
		newUser := &model.User{
			DiscordID: userResponse.ID,
			DiscordAccessToken: oauth2Response.AccessToken,
			DiscordRefreshToken: oauth2Response.RefreshToken,
		}
		createdUser, err := s.repo.Mongo.User.Create(ctx, newUser)
		if err != nil {
			s.logger.Sugar().Errorf("failed to CREATE new user(%s) in database: %s", userResponse.ID, err.Error())
			return "", "", errInternal
		}
		user = createdUser
	}

	userUpdates := map[string]interface{}{
		"discordAccessToken": oauth2Response.AccessToken,
		"discordRefreshToken": oauth2Response.RefreshToken,
	}
	if err := s.repo.Mongo.User.UpdateByID(ctx, user.ID, userUpdates); err != nil {
		s.logger.Sugar().Errorf("failed to update user(%s): %s", user.ID, err.Error())
		return "", "", errInternal
	}

	// Response tokens
	accessTokenClaims := jwt.MapClaims{
		"id": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}
	refreshTokenClaims := jwt.MapClaims{
		"id": user.ID.String(),
		"discordId": user.DiscordID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	accessToken, refreshToken, err := lib.GenerateJWTPair(accessTokenClaims, refreshTokenClaims)
	if err != nil {
		s.logger.Sugar().Errorf("failed to generate jwt pair for user(%s): %s", userResponse.ID, err.Error())
		return "", "", errInternal
	}

	return accessToken, refreshToken, nil
}
