package handlers

import (
	"net/http"
)

func HandlerHealhtz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		panic(err)
	}
}
