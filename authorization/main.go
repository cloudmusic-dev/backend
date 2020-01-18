package authorization

import (
	"fmt"
	"github.com/cloudmusic-dev/backend/configuration"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

var Config configuration.Configuration

func CreateToken(userId uuid.UUID) (string, error) {
	claim := jwt.MapClaims{
		"userId": userId.String(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().AddDate(0, 0, 14).Unix(),
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

func ValidateRequest(r *http.Request) (bool, *uuid.UUID) {
	authorization := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		return false, nil
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	return ValidateToken(token)
}

func ValidateToken(tokenString string) (bool, *uuid.UUID) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if Config.Signing.Method == "hmac" {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(Config.Signing.Key), nil
		} else {
			return nil, fmt.Errorf("invalid signing method")
		}
	})

	if err != nil {
		log.Printf("Failed to verify token: %v", err)
		return false, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, err := uuid.Parse(claims["userId"].(string))
		if err != nil {
			log.Printf("Failed to parse user id in jwt: %v", err)
			return false, nil
		}

		return true, &id
	} else {
		return false, nil
	}
}