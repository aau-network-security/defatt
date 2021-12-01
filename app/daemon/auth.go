package daemon

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aau-network-security/defatt/store"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

const (
	USERNAME_KEY    = "un"
	SUPERUSER_KEY   = "su"
	VALID_UNTIL_KEY = "vu"
)

var (
	ErrMissingToken          = errors.New("no security token provided")
	ErrInvalidUsernameOrPass = errors.New("invalid username or password")
	ErrInvalidTokenFormat    = errors.New("invalid token format")
	ErrTokenExpired          = errors.New("token has expired")
	ErrUnknownUser           = errors.New("unknown user")
	ErrEmptyUser             = errors.New("username cannot be empty")
	ErrEmptyPasswd           = errors.New("password cannot be empty")
)

type Authenticator interface {
	TokenForUser(username, password string) (string, error)
	AuthenticateContext(context.Context) (context.Context, error)
}

type auth struct {
	us  store.UserStore
	key string
}

type us struct{}

func NewAuthenticator(us store.UserStore, key string) Authenticator {
	return &auth{
		us:  us,
		key: key,
	}
}

func (a *auth) TokenForUser(username, password string) (string, error) {
	username = strings.ToLower(username)

	if username == "" {
		return "", ErrEmptyUser
	}

	if password == "" {
		return "", ErrEmptyPasswd
	}

	u, err := a.us.GetUserByUsername(username)
	if err != nil {
		return "", ErrInvalidUsernameOrPass
	}

	if ok := u.IsCorrectPassword(password); !ok {
		return "", ErrInvalidUsernameOrPass
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		USERNAME_KEY:    u.Username,
		SUPERUSER_KEY:   u.SuperUser,
		VALID_UNTIL_KEY: time.Now().Add(31 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *auth) AuthenticateContext(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, ErrMissingToken
	}

	if len(md["token"]) == 0 {
		return ctx, ErrMissingToken
	}

	token := md["token"][0]
	if token == "" {
		return ctx, ErrMissingToken
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return ctx, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.key), nil
	})
	if err != nil {
		return ctx, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return ctx, ErrInvalidTokenFormat
	}

	username, ok := claims[USERNAME_KEY].(string)
	if !ok {
		return ctx, ErrInvalidTokenFormat
	}

	u, err := a.us.GetUserByUsername(username)
	if err != nil {
		return ctx, ErrUnknownUser
	}

	validUntil, ok := claims[VALID_UNTIL_KEY].(float64)
	if !ok {
		return ctx, ErrInvalidTokenFormat
	}

	if int64(validUntil) < time.Now().Unix() {
		return ctx, ErrTokenExpired
	}

	ctx = context.WithValue(ctx, us{}, u)

	return ctx, nil
}
