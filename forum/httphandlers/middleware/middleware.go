package middleware

import (
	"net/http"
	"nx_trainee_forum/forum/application/config"
	"nx_trainee_forum/forum/httphandlers/authorization"

	"gorm.io/gorm"
)

func Authorization(cfg *config.Config, db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		u := authorization.GetCurrentUser(cfg, db, r)
		if u.ID == 0 {
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			w.Write([]byte(`{"error":""}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
