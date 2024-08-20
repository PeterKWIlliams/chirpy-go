package handlers

import (
	"fmt"
	"net/http"
)

func (cfg *ApiCfg) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++
		fmt.Printf("server hit: %d\n", cfg.FileserverHits)
		next.ServeHTTP(w, r)
	})
}
