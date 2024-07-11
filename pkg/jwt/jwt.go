package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewToken(userId int) (string, error) {
	now := time.Now().Local()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 7)),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserId: userId,
	})
	return claims.SignedString([]byte(""))
}

func ParseToken(token string) (int, error) {
	if token == "" {
		return 0, errors.New("token为空")
	}
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, errors.New("unexpected method")
		}
		return []byte(""), nil
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if !t.Valid {
		return 0, errors.New("token已失效")
	}
	if c, ok := t.Claims.(*UserClaims); ok {
		return c.UserId, nil
	}
	return 0, errors.New("token已失效")
}
