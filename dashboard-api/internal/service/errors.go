package service

import "errors"

var (
	ErrInternal = errors.New("internal server error")
	ErrTokenIsNotValid = errors.New("token is not valid")
	ErrDiscordIDRequiredIfCacheEnabled = errors.New("discordID is required if cache enabled")
	ErrIDIsNotValid = errors.New("ID is not valid")
)
