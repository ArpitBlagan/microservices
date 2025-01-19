package middleware

import (
	"context"
	"fmt"
	"go-backend/utils"
	"net/http"
)
type contextKey string
const userKey contextKey = "user"

func ValidateMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("trying to validate token")
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hehe, err := utils.ValidateJwt(cookie.Value)
		ctx:=context.WithValue(r.Context(),userKey,hehe.UserID)
		if err != nil  {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("Token is correct :)",hehe)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}