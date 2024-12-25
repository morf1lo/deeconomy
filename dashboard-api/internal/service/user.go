package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/morf1lo/deeconomy-bot-api/internal/dto"
	"github.com/morf1lo/deeconomy-bot-api/internal/lib"
	"github.com/morf1lo/deeconomy-bot-api/internal/model"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository"
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
	reqBody.Add("grant_type", "authorization_token")
	reqBody.Add("redirect_uri", os.Getenv("DISCORD_REDIRECT_URI"))

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

	fmt.Println(oauth2Resp)

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

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", oauth2Response.AccessToken)

	userResp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Sugar().Errorf("failed to DO request to get user info: %s", err.Error())
		return "", "", err
	}
	defer userResp.Body.Close()

	var userResponse dto.DiscordUserResp
	if err := json.NewDecoder(userResp.Body).Decode(&userResponse); err != nil {
		s.logger.Sugar().Errorf("failed to decode discord user response body: %s", err.Error())
		return "", "", errInternal
	}

	if userResp.StatusCode != http.StatusOK {
		s.logger.Sugar().Errorf("error from discord api(status: %d): %s", userResp.StatusCode, userResponse.Message)
		return "", "", fmt.Errorf("discord error: %s", userResponse.Message)
	}

	user := &model.User{
		DiscordID: userResponse.User.ID,
		DiscordAccessToken: oauth2Response.AccessToken,
		DiscordRefreshToken: oauth2Response.RefreshToken,
	}
	createdUser, err := s.repo.Mongo.User.Create(ctx, user)
	if err != nil {
		s.logger.Sugar().Errorf("failed to CREATE new user(%s) in database: %s", userResponse.User.ID, err.Error())
		return "", "", errInternal
	}

	accessTokenClaims := jwt.MapClaims{
		"id": createdUser.ID.String(),
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}
	refreshTokenClaims := jwt.MapClaims{
		"id": createdUser.ID.String(),
		"discordId": createdUser.DiscordID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	accessToken, refreshToken, err := lib.GenerateJWTPair(accessTokenClaims, refreshTokenClaims)
	if err != nil {
		s.logger.Sugar().Errorf("failed to generate jwt pair for user(%s): %s", userResponse.User.ID, err.Error())
		return "", "", errInternal
	}

	return accessToken, refreshToken, nil
}
