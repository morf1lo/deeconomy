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
	"github.com/morf1lo/deeconomy-bot-api/internal/repository/redisrepo"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return "", "", ErrInternal
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
		return "", "", ErrInternal
	}

	if oauth2Resp.StatusCode != http.StatusOK {
		s.logger.Sugar().Errorf("error from discord api(status: %d): %s", oauth2Resp.StatusCode, oauth2Response.Message)
		return "", "", fmt.Errorf("discord error: %s", oauth2Response.Message)
	}

	userResponse, err := s.FindDiscordUser(ctx, "", oauth2Response.AccessToken, false)
	if err != nil {
		s.logger.Sugar().Errorf("failed to fetch user from discord: %s", err.Error())
		return "", "", err
	}

	// Creating or updating user
	user, err := s.repo.Mongo.User.FindByDiscordID(ctx, userResponse.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		s.logger.Sugar().Errorf("failed to find user(%s): %s", userResponse.ID)
		return "", "", ErrInternal
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
			return "", "", ErrInternal
		}
		user = createdUser
	}

	user.DiscordAccessToken = oauth2Response.AccessToken
	user.DiscordRefreshToken = oauth2Response.RefreshToken

	userUpdates := map[string]interface{}{
		"discordAccessToken": oauth2Response.AccessToken,
		"discordRefreshToken": oauth2Response.RefreshToken,
	}
	if err := s.repo.Mongo.User.UpdateByID(ctx, user.ID, userUpdates); err != nil {
		s.logger.Sugar().Errorf("failed to update user(%s): %s", user.ID, err.Error())
		return "", "", ErrInternal
	}

	// Response tokens
	accessTokenClaims := jwt.MapClaims{
		"id": user.ID.Hex(),
		"discordId": user.DiscordID,
		"discordAccessToken": user.DiscordAccessToken,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}
	refreshTokenClaims := jwt.MapClaims{
		"id": user.ID.Hex(),
		"discordId": user.DiscordID,
		"discordRefreshToken": user.DiscordRefreshToken,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	accessToken, refreshToken, err := lib.GenerateJWTPair(accessTokenClaims, refreshTokenClaims)
	if err != nil {
		s.logger.Sugar().Errorf("failed to generate jwt pair for user(%s): %s", userResponse.ID, err.Error())
		return "", "", ErrInternal
	}

	return accessToken, refreshToken, nil
}

func (s *userService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	decodedToken, err := lib.DecodeRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrTokenIsNotValid
	}
	fmt.Println(decodedToken["discordRefreshToken"].(string))

	revokeEndpoint := fmt.Sprintf("%s/oauth2/token/revoke", DISCORD_HOST)

	reqBody := url.Values{}
	reqBody.Add("client_id", os.Getenv("CLIENT_ID"))
	reqBody.Add("client_secret", os.Getenv("CLIENT_SECRET"))
	reqBody.Add("grant_type", "refresh_token")
	reqBody.Add("refresh_token", decodedToken["discordRefreshToken"].(string))

	req, err := http.NewRequest(http.MethodPost, revokeEndpoint, strings.NewReader(reqBody.Encode()))
	if err != nil {
		s.logger.Sugar().Errorf("failed to CREATE new request to REFRESH tokens: %s", err.Error())
		return "", "", ErrInternal
	}

	revokeTokensResp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Sugar().Errorf("failed to DO request to refresh user tokens: %s", err.Error())
		return "", "", err
	}
	defer revokeTokensResp.Body.Close()

	var revokeTokensResponse dto.DiscordRevokeResp
	if err := json.NewDecoder(revokeTokensResp.Body).Decode(&revokeTokensResponse); err != nil {
		s.logger.Sugar().Errorf("failed to decode response body: %s", err.Error())
		return "", "", ErrInternal
	}

	if revokeTokensResp.StatusCode != http.StatusOK {
		s.logger.Sugar().Errorf("error from discord api(status: %d): %s", revokeTokensResp.StatusCode, revokeTokensResponse.Message)
		return "", "", fmt.Errorf("discord error: %s", revokeTokensResponse.Message)
	}

	userUpdates := map[string]interface{}{
		"discordAccessToken": revokeTokensResponse.AccessToken,
		"discordRefreshToken": revokeTokensResponse.RefreshToken,
	}
	if err := s.repo.Mongo.User.UpdateByID(ctx, decodedToken["id"].(primitive.ObjectID), userUpdates); err != nil {
		s.logger.Sugar().Errorf("failed to update user(%s) in MongoDB: %s", decodedToken["id"].(string), err.Error())
		return "", "", ErrInternal
	}

	objectID, err := primitive.ObjectIDFromHex(decodedToken["id"].(string))
	if err != nil {
		return "", "", ErrIDIsNotValid
	}

	accessTokenClaims := jwt.MapClaims{
		"id": objectID.Hex(),
		"discordId": decodedToken["discordId"].(string),
		"discordAccessToken": revokeTokensResponse.AccessToken,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}
	refreshTokenClaims := jwt.MapClaims{
		"id": objectID.Hex(),
		"discordId": decodedToken["discordId"].(string),
		"discordRefreshToken": revokeTokensResponse.RefreshToken,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	accessToken, refreshToken, err := lib.GenerateJWTPair(accessTokenClaims, refreshTokenClaims)
	if err != nil {
		s.logger.Sugar().Errorf("failed to generate jwt pair: %s", err.Error())
		return "", "", ErrInternal
	}

	return accessToken, refreshToken, nil
}

func (s *userService) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	userCache, err := redisrepo.Get[model.User](s.repo.Redis.Default, ctx, redisrepo.UserKey(id.Hex()))
	if err == nil {
		return userCache, nil
	}
	if err != redis.Nil {
		s.logger.Sugar().Errorf("failed to get user(%s) from Redis: %s", id.String(), err.Error())
		return nil, ErrInternal
	}

	user, err := s.repo.Mongo.User.FindByID(ctx, id)
	if err != nil {
		s.logger.Sugar().Errorf("failed to find user(%s) in MongoDB: %s", id.Hex(), err.Error())
		return nil, err
	}

	if err := s.repo.Redis.Default.SetJSON(ctx, redisrepo.UserKey(id.Hex()), user, time.Hour * 3); err != nil {
		s.logger.Sugar().Errorf("failed to set user(%s) in Redis: %s", id.Hex(), err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userService) FindDiscordUser(ctx context.Context, discordID string, discordAccessToken string, cacheEnable bool) (*discordgo.User, error) {
	discordID = strings.TrimSpace(discordID)
	if discordID == "" {
		cacheEnable = false
	}

	if cacheEnable {
		discordUser, err := redisrepo.Get[discordgo.User](s.repo.Redis.Default, ctx, redisrepo.DiscordUserKey(discordID))
		if err == nil {
			return discordUser, nil
		}
		if err != redis.Nil {
			s.logger.Sugar().Errorf("failed to get discord user(%s) from Redis: %s", discordUser, err.Error())
			return nil, err
		}
	}

	usersEndpoint := fmt.Sprintf("%s/users/@me", DISCORD_HOST)
	req, err := http.NewRequest(http.MethodGet, usersEndpoint, nil)
	if err != nil {
		s.logger.Sugar().Errorf("failed to CREATE request for '/users/@me' discord endpoint: %s", err.Error())
		return nil, ErrInternal
	}

	req.Header.Add("Authorization", "Bearer " + discordAccessToken)

	userResp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Sugar().Errorf("failed to DO request to get user info: %s", err.Error())
		return nil, err
	}
	defer userResp.Body.Close()

	userRespBodyBytes, err := io.ReadAll(userResp.Body)
	if err != nil {
		s.logger.Sugar().Errorf("failed to read discord user response body: %s", err.Error())
		return nil, ErrInternal
	}

	if userResp.StatusCode != http.StatusOK {
		var discordError dto.DiscordError
		if err := json.Unmarshal(userRespBodyBytes, &discordError); err != nil {
			s.logger.Sugar().Errorf("failed to decode discord error response: %s", err.Error())
			return nil, fmt.Errorf("discord error: unknown error")
		}

		s.logger.Sugar().Errorf("error from discord api(status: %d, code: %d): %s", userResp.StatusCode, discordError.Code, discordError.Message)
		return nil, fmt.Errorf("discord error: %s", discordError.Message)
	}

	var userResponse discordgo.User
	if err := json.Unmarshal(userRespBodyBytes, &userResponse); err != nil {
		s.logger.Sugar().Errorf("failed to decode discord user response: %s", err.Error())
		return nil, ErrInternal
	}

	if cacheEnable {
		if err := s.repo.Redis.Default.SetJSON(ctx, redisrepo.DiscordUserKey(discordID), &userResponse, time.Minute * 30); err != nil {
			s.logger.Sugar().Errorf("failed to set discord user(%s) in Redis: %s", discordID, err.Error())
			return nil, err
		}
	}

	return &userResponse, nil
}
