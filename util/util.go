package util

import (
	"net/http"
	"fmt"
)

func Must(err error){
	if err != nil {
		panic(err)
	}
}

func HTTPBadRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request:", err)

	return true
}
