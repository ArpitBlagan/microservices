package middleware

import (
	"fmt"
	"go-backend/utils"
	"net/http"
)

func ValidateMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hehe, err := utils.ValidateJwt(cookie.Value)

		if err != nil  {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("Token is correct :)",hehe)
		next.ServeHTTP(w, r)
	})
}