package api

import (
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	AccessToken          string    `json:"access_token"`
}
