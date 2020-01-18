package authorization

import (
	"fmt"
	"github.com/cloudmusic-dev/backend/configuration"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var Config configuration.Configuration

func CreateToken(userId uuid.UUID) (string, error) {
	claim := jwt.MapClaims{
		"userId": userId.String(),
	}

	if Config.Signing.Method == "hmac" {
		var signing *jwt.SigningMethodHMAC
		if Config.Signing.Strength == 256 {
			signing = jwt.SigningMethodHS256
		} else if Config.Signing.Strength == 384 {
			signing = jwt.SigningMethodHS384
		} else if Config.Signing.Strength == 512 {
			signing = jwt.SigningMethodHS512
		} else {
			return "", fmt.Errorf("invalid signing strength for hmac, valid values: 256, 384, 512")
		}

		token := jwt.NewWithClaims(signing, claim)
		return token.SignedString([]byte(Config.Signing.Key))
	}

	return "", fmt.Errorf("invalid signing method")
}