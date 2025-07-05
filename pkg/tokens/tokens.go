package tokens

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptySigningKey  = errors.New("empty signing key when create tokens.Manager in New")
	ErrSignedAccess     = errors.New("error when signing token in NewAccessToken")
	ErrAccessParsing    = errors.New("error when parsing access token in ParseAccessToken")
	ErrTypeAssertion    = errors.New("error when type assertion claims in ParseAccessToken")
	ErrSubField         = errors.New("error when GetSubject from claims in ParseAccessToken")
	ErrGenerateRefreshT = errors.New("error when rand.Read in NewRefreshToken")
	ErrSignedRefresh    = errors.New("error when signing token in NewRefreshToken")
	ErrRefreshParsing   = errors.New("error when parsing refresh token in ParseRefreshToken")
	ErrDecodeRefresh    = errors.New("error when decode base64 refresh token in ParseRefreshToken")
)

type Manager interface {
	NewAccessToken(userId string) (string, error)
	ParseAccessToken(accessTkn string) (string, error)
	NewRefreshToken() (string, error)
	ParseRefreshToken(refreshTkn string) (string, error)
}

type managerStrct struct {
	signingKey    string
	accessTknTTL  time.Duration
	refreshTknTTL time.Duration
}

func New(signingKey string, accessTknTTL, refreshTknTTL time.Duration) (Manager, error) {
	if signingKey == "" {
		return nil, ErrEmptySigningKey
	}

	return &managerStrct{
		signingKey,
		accessTknTTL,
		refreshTknTTL,
	}, nil
}

// Создание jwt токена доступа. Стандартные claims, решил, что кастомные будут излишком
func (m *managerStrct) NewAccessToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(m.accessTknTTL).Unix(),
	})

	signedTkn, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", errors.Join(ErrSignedAccess, err)
	}

	return signedTkn, nil
}

// Парсинг jwt токена доступа. Вернет содержимое поле sub из claims.
func (m *managerStrct) ParseAccessToken(accessTkn string) (string, error) {
	token, err := jwt.Parse(accessTkn, func(t *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		return "", errors.Join(ErrAccessParsing, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId, err := claims.GetSubject()
		if err != nil {
			return "", errors.Join(ErrSubField, err)
		}
		return userId, nil
	} else {
		return "", errors.Join(ErrTypeAssertion, err)
	}
}

// Создание jwt рефреш токена. Вернет подписанный токен в base64
func (m *managerStrct) NewRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(m.refreshTknTTL).Unix(),
	})

	signedTkn, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", errors.Join(ErrSignedRefresh, err)
	}
	tokenStr := base64.URLEncoding.EncodeToString([]byte(signedTkn))
	return tokenStr, nil
}

// Парсинг jwt рефреш токена формата base64.
// Приведет к стандартносу виду jwt, проверит и вернет токен.
func (m *managerStrct) ParseRefreshToken(refreshTknBase64 string) (string, error) {
	tokenStr, err := base64.URLEncoding.DecodeString(refreshTknBase64)
	if err != nil {
		return "", errors.Join(ErrDecodeRefresh, err)
	}
	token, err := jwt.Parse(string(tokenStr), func(t *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", errors.Join(ErrRefreshParsing, err)
	}
	return token.Raw, nil
}
