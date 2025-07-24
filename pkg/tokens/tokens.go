package tokens

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptySigningKey     = errors.New("error empty signing key when create tokens.Manager in New")
	ErrSignedAccess        = errors.New("error when signing token in NewAccessToken")
	ErrAccessParsing       = errors.New("error when parsing access token in ParseAccessToken")
	ErrAccessTokenExpired  = errors.New("error access token expired")
	ErrSignedRefresh       = errors.New("error when signing token in NewRefreshToken")
	ErrRefreshParsing      = errors.New("error when parsing refresh token in ParseRefreshToken")
	ErrRefreshTokenExpired = errors.New("error refresh token expired")
)

type Manager interface {
	// Создание jwt токена доступа. С кастомными claims
	NewAccessToken(claims jwt.Claims) (string, error)

	// Парсинг jwt токена доступа.
	// В claims нужно передать ссылку на структуру.
	// Вернет ссылку на jwt токен.
	// Если токен просрочен, то вренет ссылку на токен и ошибку.
	ParseAccessToken(accessTkn string, claims jwt.Claims) (*jwt.Token, error)

	// Создание jwt рефреш токена. С кастомными claims
	NewRefreshToken(claims jwt.Claims) (string, error)

	// Парсинг jwt рефреш токена.
	// В claims нужно передать ссылку на структуру.
	// Вернет ссылку на jwt токен.
	ParseRefreshToken(refreshTkn string, claims jwt.Claims) (*jwt.Token, error)

	// Объединение создания токена доступа и рефреш токена.
	// Вернет структур с парой токенов.
	NewTokensPair(accessClaims, refreshClaims jwt.Claims) (TokensPair, error)

	// Парсинг jwt токена доступа.
	// В claims нужно передать ссылку на структуру.
	// Вернет структуру с ссылками на jwt токены.
	ParseTokenPair(tokenA, tokenR string, claimsA, claimsR jwt.Claims) (ParsedTokensPair, error)
}

type TokensPair struct {
	AccessToken  string
	RefreshToken string
}

type ParsedTokensPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}

type manager struct {
	signingKey string
}

func New(signingKey string) (Manager, error) {
	if signingKey == "" {
		return nil, ErrEmptySigningKey
	}

	return &manager{
		signingKey,
	}, nil
}

func (m *manager) NewAccessToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedTkn, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", errors.Join(ErrSignedAccess, err)
	}

	return signedTkn, nil
}

func (m *manager) ParseAccessToken(token string, claims jwt.Claims) (*jwt.Token, error) {
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return tkn, errors.Join(ErrAccessTokenExpired, err)
		}
		return &jwt.Token{}, errors.Join(ErrAccessParsing, err)
	}

	return tkn, nil
}

func (m *manager) NewRefreshToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedTkn, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", errors.Join(ErrSignedRefresh, err)
	}

	return signedTkn, nil
}

func (m *manager) ParseRefreshToken(token string, claims jwt.Claims) (*jwt.Token, error) {
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	}, jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return &jwt.Token{}, errors.Join(ErrRefreshTokenExpired, err)
		}
		return &jwt.Token{}, errors.Join(ErrRefreshParsing, err)
	}
	return tkn, nil
}

func (m *manager) NewTokensPair(claimsA, claimsR jwt.Claims) (TokensPair, error) {
	accessT, err := m.NewAccessToken(claimsA)
	if err != nil {
		return TokensPair{}, err
	}

	refreshT, err := m.NewRefreshToken(claimsR)
	if err != nil {
		return TokensPair{}, err
	}

	return TokensPair{accessT, refreshT}, nil
}

func (m *manager) ParseTokenPair(
	tokenA, tokenR string,
	claimsA, claimsR jwt.Claims,
) (ParsedTokensPair, error) {
	accessT, err := m.ParseAccessToken(tokenA, claimsA)
	if err != nil {
		return ParsedTokensPair{}, err
	}

	refreshT, err := m.ParseRefreshToken(tokenR, claimsR)
	if err != nil {
		return ParsedTokensPair{}, err
	}

	return ParsedTokensPair{accessT, refreshT}, nil
}
