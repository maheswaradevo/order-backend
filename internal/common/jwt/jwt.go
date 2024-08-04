package jwt

import (
	"fmt"
	"order-service-backend/internal/config"
	"time"

	"github.com/golang-jwt/jwt"
)

func getJwtSecretKey() []byte {
	jwtSecretKey := config.GetConfig().JWTConfig.SecretKey
	return []byte(jwtSecretKey)
}

func Encode(data map[string]interface{}) (result string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	for key, value := range data {
		claims[key] = value
	}

	accessToken, errSignString := token.SignedString(getJwtSecretKey())
	if errSignString != nil {
		return
	}

	result = accessToken
	return
}

func Decode(tokenString string) (result map[string]interface{}, err error) {
	token, errParseToken := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[error] Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return getJwtSecretKey(), nil
	})
	if errParseToken != nil {
		err = errParseToken
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result = claims
		return
	}

	return
}

func createAccessToken(data map[string]interface{}) (*string, *time.Duration, error) {
	jwtTimeout, _ := time.ParseDuration(config.GetConfig().JWTConfig.Timeout)
	exp := time.Now().Add(jwtTimeout).Unix()

	data["exp"] = exp

	accessToken, err := Encode(data)
	if err != nil {
		return nil, nil, err
	}
	return &accessToken, &jwtTimeout, err
}

func createRefreshToken(data map[string]interface{}) (*string, *time.Duration, error) {
	jwtTimeout, _ := time.ParseDuration(config.GetConfig().JWTConfig.Timeout)
	exp := time.Now().Add(jwtTimeout * 2).Unix()

	data["exp"] = exp

	refreshToken, err := Encode(data)
	if err != nil {
		return nil, nil, err
	}
	return &refreshToken, &jwtTimeout, err
}

func GenerateTokenPair(data map[string]interface{}) (*string, *string, *time.Duration, *time.Duration, error) {
	accessToken, expiredToken, err := createAccessToken(data)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	refreshToken, expRefreshToken, err := createRefreshToken(data)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return accessToken, refreshToken, expiredToken, expRefreshToken, nil
}
