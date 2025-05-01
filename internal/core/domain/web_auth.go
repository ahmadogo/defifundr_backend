package domain

import (
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

// Define the claims struct first, similar to the working example
type Web3AuthClaims struct {
	jwtv4.StandardClaims
	Email             string   `json:"email"`
	Name              string   `json:"name"`
	ProfileImage      string   `json:"profileImage"`
	Verifier          string   `json:"verifier"`
	VerifierID        string   `json:"verifierId"`
	AggregateVerifier string   `json:"aggregateVerifier"`
	Wallets           []Wallet `json:"wallets"`
	Nonce             string   `json:"nonce"`
}

// Wallet represents a wallet in Web3Auth tokens
type Wallet struct {
	PublicKey string `json:"public_key"`
	Type      string `json:"type"`
	Curve     string `json:"curve"`
}
