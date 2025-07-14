package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptySigningKey  = errors.New("error empty signing key when create tokens.Manager in New")
	ErrSignedAccess     = errors.New("error when signing token in NewAccessToken")
	ErrAccessParsing    = errors.New("error when parsing access token in ParseAccessToken")
	ErrTypeAssertion    = errors.New("error when type assertion claims in ParseAccessToken")
	ErrSubField         = errors.New("error when GetSubject from claims in ParseAccessToken")
	ErrGenerateRefreshT = errors.New("error when rand.Read in NewRefreshToken")
	ErrSignedRefresh    = errors.New("error when signing token in NewRefreshToken")
	ErrRefreshParsing   = errors.New("error when parsing refresh token in ParseRefreshToken")
	ErrTokenExpired     = errors.New("error token expired")
)

type Manager interface {
	NewAccessToken(userId string) (string, error)
	ParseAccessToken(accessTkn string) (string, error)
	NewRefreshToken() (string, error)
	ParseRefreshToken(refreshTkn string) (string, error)
	NewTokensPair(userId string) (TokensPair, error)
}

type TokensPair struct {
	AccessToken     string
	RefreshToken    string
	RefreshTokenTtl time.Duration
}

type manager struct {
	signingKey    string
	accessTknTTL  time.Duration
	refreshTknTTL time.Duration
}

func New(signingKey string, accessTknTTL, refreshTknTTL time.Duration) (Manager, error) {
	if signingKey == "" {
		return nil, ErrEmptySigningKey
	}

	return &manager{
		signingKey,
		accessTknTTL,
		refreshTknTTL,
	}, nil
}

// Создание jwt токена доступа. Стандартные claims, решил, что кастомные будут излишком
func (m *manager) NewAccessToken(userId string) (string, error) {
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
func (m *manager) ParseAccessToken(accessTkn string) (string, error) {
	token, err := jwt.Parse(accessTkn, func(t *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", errors.Join(ErrTokenExpired, err)
		}
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

// Создание jwt рефреш токена. Вернет подписанный токен
func (m *manager) NewRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(m.refreshTknTTL).Unix(),
	})

	signedTkn, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", errors.Join(ErrSignedRefresh, err)
	}

	return signedTkn, nil
}

// Парсинг jwt рефреш токена. Проверит и вернет токен.
func (m *manager) ParseRefreshToken(refreshTkn string) (string, error) {
	token, err := jwt.Parse(refreshTkn, func(t *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", errors.Join(ErrRefreshParsing, err)
	}
	return token.Raw, nil
}

// Объединение создания токена доступа и рефреш токена.
// Вернет структур с парой токенов.
func (m *manager) NewTokensPair(userId string) (TokensPair, error) {
	accessT, err := m.NewAccessToken(userId)
	if err != nil {
		return TokensPair{}, err
	}

	refreshT, err := m.NewRefreshToken()
	if err != nil {
		return TokensPair{}, err
	}

	return TokensPair{accessT, refreshT, m.refreshTknTTL}, nil
}
