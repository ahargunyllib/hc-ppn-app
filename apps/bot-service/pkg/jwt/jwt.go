package jwt

import (
	"fmt"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/internal/infra/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock/mock_jwt.go -package=mock github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/jwt CustomJwtInterface

type CustomJwtInterface interface {
	Create() (string, error)
	Decode(tokenString string, claims *Claims) error
}

type Claims struct {
	jwt.RegisteredClaims
}

type CustomJwtStruct struct {
	SecretKey   string
	ExpiredTime time.Duration
}

var Jwt = getJwt()

func getJwt() CustomJwtInterface {
	return &CustomJwtStruct{
		SecretKey:   env.AppEnv.JwtSecretKey,
		ExpiredTime: env.AppEnv.JwtExpTime,
	}
}

func (j *CustomJwtStruct) Create() (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "hc-ppn-app",
			Audience:  jwt.ClaimStrings{"hc-ppn-app"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiredTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func (j *CustomJwtStruct) Decode(tokenString string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}
