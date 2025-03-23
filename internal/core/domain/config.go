package domain

import "time"

// Config holds application configuration
type Config struct {
	// JWT Configuration
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	TokenSymmetricKey    string
	HTTPServerAddress   string
	
	// Database Configuration
	DBDriver   string
	DBSource   string
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	
	// Server Configuration
	ServerAddress string
	Environment   string
	
	// Email Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	
	// Security Configuration
	RotateRefreshTokens bool
	MaxLoginAttempts    int
	
	// OTP Configuration
	OTPExpiryDuration time.Duration
	MaxOTPAttempts     int
}