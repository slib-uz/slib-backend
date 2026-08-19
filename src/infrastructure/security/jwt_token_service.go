package security

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/security"
	"slib.uz/src/infrastructure/config"
)

type JwtTokenService struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// @inject
func NewJwtTokenService(cfg *config.Config) security.TokenService {
	return &JwtTokenService{
		privateKey: loadPrivateKey(cfg.JwtPrivateKeyPath),
		publicKey:  loadPublicKey(cfg.JwtPublicKeyPath),
	}
}

func (s *JwtTokenService) Encode(data *entity.TokenEntity) string {

	claim := NewJwtTokenClaim(
		jwt.RegisteredClaims{
			Subject:   data.Subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(data.Exp),
			ID:        uuid.NewString(),
		},
		data.Payload,
	)

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)

	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return ""
	}

	return signedToken
}

func (s *JwtTokenService) Decode(tokenString string) (*entity.TokenEntity, error) {

	token, err := jwt.ParseWithClaims(tokenString, &JwtTokenClaim{}, func(token *jwt.Token) (any, error) {

		// 🔒 ALG spoofingdan himoya
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, response.InvalidTokenError
		}

		return s.publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, response.ExpiredTokenError
		}
		return nil, response.InvalidTokenError
	}

	claims, ok := token.Claims.(*JwtTokenClaim)
	if !ok || !token.Valid {
		return nil, response.InvalidTokenError
	}

	decoded := entity.NewTokenEntity(
		claims.ExpiresAt.Time,
		claims.Subject,
		claims.Payload,
	)
	decoded.ID = claims.ID

	return decoded, nil
}

func loadPrivateKey(path string) *ecdsa.PrivateKey {
	key, _ := os.ReadFile(path)
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(key)
	if err != nil {
		panic(fmt.Errorf("[JwtTokenService] failed to load private key: %w", err))
	}
	return privateKey
}

func loadPublicKey(path string) *ecdsa.PublicKey {
	key, _ := os.ReadFile(path)
	publicKey, err := jwt.ParseECPublicKeyFromPEM(key)
	if err != nil {
		panic(fmt.Errorf("[JwtTokenService] failed to load public key: %w", err))
	}
	return publicKey
}
