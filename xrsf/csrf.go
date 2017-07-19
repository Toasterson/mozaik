package xrsf

import (
	"net/http"
	"github.com/justinas/nosurf"
	"github.com/toasterson/mozaik/logger"
)

func Nosurfing(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(fail))
	return surfing
}


func fail(w http.ResponseWriter, r *http.Request) {
	logger.Warn("Failed to validate XSRF Token:", nosurf.Reason(r))
	w.WriteHeader(http.StatusBadRequest)
}