package interfaces

import "time"

// Request struct for creating a wallet
type CreateWalletAddressRequest struct {
	WalletAddress string `json:"wallet_address" binding:"required"`
	Chain         string `json:"chain" binding:"required"`
}

type CreateWalletAddressResponse struct {
	ID            int64     `json:"id"`
	UserID        string    `json:"user_id"`
	WalletAddress string    `json:"wallet_address"`
	Chain         string    `json:"chain"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}