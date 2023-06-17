package interfaces

import (
	"mime/multipart"
	"time"
)

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type GetPrivateKeyRequest struct {
	Password string `json:"password" binding:"required"`
}

type AddressResponse struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

type CreateCampaignTypeRequest struct {
	Name  string          `json:"name" binding:"required"`
	Image *multipart.File `json:"image" binding:"required"`
}
