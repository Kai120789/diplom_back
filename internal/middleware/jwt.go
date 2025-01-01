package middleware

import (
	"context"
	"fmt"
	"materials/internal/config"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDContextKey   contextKey = "userID"
	UserNameContextKey contextKey = "userName"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("accesstoken")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		secretKey := []byte(config.AppConfig.JWTSecret)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var userID int
			var okID bool
			if IDVal, ok := claims["userID"]; ok {
				switch v := IDVal.(type) {
				case string:
					userID, err = strconv.Atoi(v)
					if err != nil {
						http.Error(w, "unauthorized", http.StatusUnauthorized)
						return
					}
					okID = true
				case float64:
					userID = int(v)
					okID = true
				}
			}

			username, okName := claims["userName"].(string)
			if !okID || !okName {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			ctx = context.WithValue(ctx, UserNameContextKey, username)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
