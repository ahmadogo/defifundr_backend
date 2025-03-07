package domain

import (
	"time"

	"github.com/google/uuid"
)

type BlockchainType string

const (
	BlockchainEthereum BlockchainType = "ethereum"
	BlockchainSolana   BlockchainType = "solana"
)

type Wallet struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"user_id"`
	WalletAddress string         `json:"wallet_address"`
	Chain         BlockchainType `json:"chain"`
	IsPrimary     bool           `json:"is_primary"`
	CreatedAt     time.Time      `json:"created_at"`
}

type CreateWalletParams struct {
	UserID        uuid.UUID      `json:"user_id" validate:"required"`
	WalletAddress string         `json:"wallet_address" validate:"required"`
	Chain         BlockchainType `json:"chain" validate:"required"`
	IsPrimary     bool           `json:"is_primary"`
}

type UpdateWalletPrimaryParams struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	IsPrimary bool      `json:"is_primary"`
}