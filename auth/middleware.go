package auth

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/toasterson/mozaik/logger"
)

type authProtector struct {
	f http.HandlerFunc
}

func AuthProtect(f http.HandlerFunc) authProtector {
	return authProtector{f}
}

func (ap authProtector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if u, err := Ab.CurrentUser(w, r); err != nil {
		logger.Error("Error fetching current user:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if u == nil {
		logger.Trace("Redirecting unauthorized user from:", r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		ap.f(w, r)
	}
}

func Nosurfing(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Warn("Failed to validate XSRF Token:", nosurf.Reason(r))
		w.WriteHeader(http.StatusBadRequest)
	}))
	return surfing
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Trace(r.Method, r.URL.Path, r.Proto)
		session, err := sessionStore.Get(r, sessionCookieName)
		if err == nil {
			logger.Trace("Session: ")
			first := true
			for k, v := range session.Values {
				if first {
					first = false
				} else {
					logger.Trace(", ")
				}
				logger.Trace(k, v)
			}
		}
		logger.Trace("Database: ")
		for _, u := range database.Users {
			logger.Trace(u)
		}
		h.ServeHTTP(w, r)
	})
}
