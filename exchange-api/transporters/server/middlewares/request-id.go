package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const requestIdHeader = "x-request-id"
const requestIdContextKey = "requestId"

func AddRequestId(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestId string

		if requestId = r.Header.Get(requestIdHeader); requestId == "" {
			requestId = uuid.New().String()
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, requestIdContextKey, requestId)

		request := r.WithContext(ctx)

		handler.ServeHTTP(w, request)
	})
}
