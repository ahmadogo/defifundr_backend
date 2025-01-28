package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	CryptoDeployURL      string        `mapstructure:"CRYPT_DEPLOY_URL"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	ContractPrivateKey   string        `mapstructure:"CONTRACT_PRIVATE_KEY"`
	ContractAddress      string        `mapstructure:"CONTRACT_ADDRESS"`
	CloudinaryURL        string        `mapstructure:"CLOUDINARY_API_KEY"`
	DeployKey            string        `mapstructure:"DEPLOY_PRIVATE_KEY"`
	DeployAddress        string        `mapstructure:"DEPLOY_ADDR"`
	Email                string        `mapstructure:"EMAIL"`
	EmailPass            string        `mapstructure:"EMAIL_PASS"`
	PassPhase            string        `mapstructure:"PASS_PHASE"`
	RedisHost            string        `mapstructure:"REDIS_HOST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName(".env")
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.AllowEmptyEnv(true)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}

// convert timestamp to unix for solidity
func ConvertToUnix(timestamp time.Time) int64 {
	return timestamp.Unix()
}

// convert unix to timestamp golang
func ConvertToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}
