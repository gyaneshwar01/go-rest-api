package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func WithJWTAuth(handler http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)

		if err != nil {
			log.Printf("failed to authenticate token")
			WriteJson(w, http.StatusUnauthorized, ErrorResponse{Error: fmt.Errorf("permission denied").Error()})
			return
		}

		if !token.Valid {
			log.Printf("failed to authenticate token")
			WriteJson(w, http.StatusUnauthorized, ErrorResponse{Error: fmt.Errorf("permission denied").Error()})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["userID"].(string)

		_, err = store.GetUserById(userId)
		if err != nil {
			log.Println("Failed to get user")
			WriteJson(w, http.StatusUnauthorized, ErrorResponse{Error: fmt.Errorf("permission denied").Error()})
			return
		}

		handler(w, r)
	}
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func validateJWT(t string) (*jwt.Token, error) {
	secret := Envs.JWTSecret

	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func CreateJWT(secret []byte, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(int(id)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
