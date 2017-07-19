package controller

import (
	"net/http"
	"github.com/toasterson/mozaik/logger"
	"errors"
)

type Handler struct {
	ControllerInterface
}

func (this *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this.IsAuthProtected() {
		if u, err := this.GetAuth().CurrentUser(w, r); err != nil {
			logger.Error("Error fetching current user:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if u == nil {
			logger.Trace("Redirecting unauthorized user from:", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	this.Prepare(w, r)
	defer this.Finish()
	var err error = nil
	switch r.Method {
	case http.MethodGet:
		err = this.Get()
	case http.MethodPost:
		err = this.Post()
	case http.MethodPatch:
		err = this.Patch()
	case http.MethodPut:
		err = this.Put()
	case http.MethodHead:
		err = this.Head()
	case http.MethodDelete:
		err = this.Delete()
	case http.MethodOptions:
		err = this.Options()
	default:
		err = errors.New("No Method")
	}
	//TODO Better Handling Should an error Occur Maybe Custom error Class?
	if err != nil {
		logger.Error(err)
		NotFound(w, r)

	}
}
