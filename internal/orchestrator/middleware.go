package orchestrator

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)


func JWTMiddleware(secretKey []byte, exemptPaths map[string]bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			// Пропускаем проверку токена для разрешённых путей
			if exemptPaths[path] || strings.HasPrefix(path, "/css/") || strings.HasPrefix(path, "/js/") {
				next.ServeHTTP(w, r)
				return
			}

			// Проверка токена
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return secretKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Добавляем userID в контекст
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userID := uint(claims["user_id"].(float64))
				ctx := context.WithValue(r.Context(), "userID", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, `{"error":"Invalid claims"}`, http.StatusUnauthorized)
			}
		})
	}
}


// Получение userID из контекста
func GetUserID(r *http.Request) (uint, bool) {
	uid, ok := r.Context().Value("userID").(uint)
	return uid, ok
}
