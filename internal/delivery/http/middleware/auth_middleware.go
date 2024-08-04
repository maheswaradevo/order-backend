package middleware

import (
	"errors"
	"strings"

	"github.com/maheswaradevo/order-backend/internal/config"

	"github.com/golang-jwt/jwt"
	middleware "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	secretKey []byte
}

func NewAuthMiddleware(secretKey []byte) *authMiddleware {
	return &authMiddleware{
		secretKey: secretKey,
	}
}

func (m *authMiddleware) AuthMiddleware() echo.MiddlewareFunc {
	return middleware.WithConfig(middleware.Config{
		SigningKey:     m.secretKey,
		SigningMethod:  middleware.AlgorithmHS256,
		ContextKey:     "userData",
		ParseTokenFunc: m.parseTokenFunc,
		TokenLookup:    "header:" + echo.HeaderAuthorization,
	})
}

func (m *authMiddleware) parseTokenFunc(c echo.Context, auth string) (interface{}, error) {
	cfg := config.GetConfig()
	headerToken := c.Request().Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		return nil, errors.New("sign in to proceed")
	}

	stringToken := strings.Split(headerToken, " ")[1]

	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTConfig.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
