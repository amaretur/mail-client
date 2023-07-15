package service

import (
	"log"
	"time"
	"errors"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/amaretur/mail-client/internal/dto"
)

var (
	ErrIncorrectTokenClaims = errors.New("incorrect token claims")
	ErrInvalidToken = errors.New("invalid token")
	ErrCreateAccessToken = errors.New("failed to create access token")
	ErrCreateRefreshToken = errors.New("failed to create refresh token")
)

type AccessClaims struct {
	Id			uint	`json:"id"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	Id			uint
	jwt.StandardClaims
}

type AuthService struct {
	secretKey		string

	accessExpires	int
	refreshExpires	int
}

func NewAuthService(
	secretKey string, accessExpires, refreshExpires int,
) *AuthService {

	return &AuthService{
		secretKey: secretKey,
		accessExpires: accessExpires,
		refreshExpires: refreshExpires,
	}
}

func (a *AuthService) CreateAccessToken(id uint) (string, error) {

	claims := AccessClaims {
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(a.accessExpires) * time.Minute,
			).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secretKey))
}

func (a *AuthService) CreateRefreshToken(id uint) (string, error) {

	claims := RefreshClaims {
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(a.refreshExpires) * time.Minute,
			).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secretKey))
}


func (a *AuthService) CreateJwtTokens(uid uint) (*dto.JwtTokens, error) {

	access, err := a.CreateAccessToken(uid)
	if err != nil {
		log.Println("access token creation error: ", err)
		return nil, ErrCreateAccessToken
	}

	refresh, err := a.CreateRefreshToken(uid)
	if err != nil {
		log.Println("refresh token creation error: ", err)
		return nil, ErrCreateRefreshToken
	}

	return &dto.JwtTokens{
		Access: access,
		Refresh: refresh,
	}, nil
}

func (a *AuthService) ParseRefreshToken(tokenStr string) (uint, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, 
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return []byte(a.secretKey), nil
		},
	)

	if err != nil {
		log.Println(err)
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return 0, ErrIncorrectTokenClaims
	}

	return claims.Id, nil
}

func (a *AuthService) ParseAccessToken(
	tokenStr string) (uint, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, 
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return []byte(a.secretKey), nil
		},
	)

	if err != nil {
		log.Println(err)
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return 0, ErrIncorrectTokenClaims
	}

	return claims.Id, nil
}




