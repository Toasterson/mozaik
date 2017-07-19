package auth

import (
	"net/http"
	"github.com/toasterson/mozaik/logger"
	"fmt"
)

//TODO Database package
var (
	Database = NewMemStorer()
)

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Trace(r.Method, r.URL.Path, r.Proto)
		session, err := sessionStore.Get(r, sessionCookieName)
		if err == nil {
			logger.Trace(fmt.Sprintf("Session: %s", session.Values))
		}
		logger.Trace(fmt.Sprintf("Database: %s", Database.Users))
		h.ServeHTTP(w, r)
	})
}
