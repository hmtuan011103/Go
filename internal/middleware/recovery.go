package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gostructure/app/pkg/response"
)

// Recovery middleware recovers from panics and returns 500 error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
				response.Error(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
