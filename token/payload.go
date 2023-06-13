package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at" `
	ExpiresAt time.Time `json:"expires"`
}

func NewPayLoad(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	} else {
		payload := &Payload{
			ID:    tokenId,
			Username:  username,
			IssuedAt:  time.Now(),
			ExpiresAt: time.Now().Add(duration),
		}
		return payload, nil
	}
}


func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
