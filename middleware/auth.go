package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/thecipherdev/goauth/utils"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenStr)

		if err != nil {
			if strings.Contains(err.Error(), "token expired") {
				refreshCookie, cookieErr := r.Cookie("refresh_token")

				if cookieErr != nil {
					http.Error(w, "Refresh token missing", http.StatusUnauthorized)
					return
				}

				refreshClaims, refreshErr := utils.ValidateRefreshToken(refreshCookie.Value)

				if refreshErr != nil {
					http.Error(w, "Refresh token invalid", http.StatusUnauthorized)
					return
				}

				newAccess, genErr := utils.GenerateToken(
					refreshClaims.Subject,
					refreshClaims.Username,
					"access",
					15*time.Minute,
				)

				if genErr != nil {
					http.Error(w, "Failed to generate new token", http.StatusUnauthorized)
					return
				}

				w.Header().Set("Authorization", "Bearer "+newAccess)
				ctx := context.WithValue(r.Context(), "userID", refreshClaims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return

			}
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
