package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	key := os.Getenv("SECRET_KEY")
	if key == "" {
		return "", errors.New("Can not find secret key")
	}

	return token.SignedString([]byte(key))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{Issuer: "chirpy"},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	claimSubject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(claimSubject)
}

func GetBearerToken(headers http.Header) (string, error) {
	a := headers["Authorization"]

	if len(a) < 1 {
		return "", errors.New("You need to provide auth token header")
	}

	token := strings.TrimPrefix(strings.Join(a, " "), "Bearer ")

	return strings.TrimSpace(token), nil
}
