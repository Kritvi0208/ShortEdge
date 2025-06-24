// F:/urlify/middleware/middleware.go
package middleware

import (
	"gofr.dev/pkg/gofr"
)

func RedirectMiddleware(h gofr.Handler) gofr.Handler {
	return func(ctx *gofr.Context) (interface{}, error) {
		// Middleware logic goes here
		return h(ctx)
	}
}
