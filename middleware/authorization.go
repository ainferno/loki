package middleware

import (
	"loki/models"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Authorization(db *gorm.DB) mux.MiddlewareFunc {
	users := models.NewUsers(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			auth := strings.Split(authHeader, " ")
			if len(auth) != 2 {
				http.Error(rw, "Bad Authorization", http.StatusUnauthorized)
				return
			}

			token := auth[1]
			user := models.User{}

			err := users.FindByToken(&user, token)
			if err != nil {
				http.Error(rw, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
