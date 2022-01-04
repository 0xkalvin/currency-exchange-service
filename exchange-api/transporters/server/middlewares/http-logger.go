package middlewares

import (
	"exchange-api/utils/logger"
	"time"

	"net/http"
)

type (
	responseWriterWrapper struct {
		http.ResponseWriter
		size       int
		statusCode int
	}
)

func (r *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size

	return size, err
}

func (r *responseWriterWrapper) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

var log = logger.NewLogger()

func HttpLogger(handler *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		wrappedResponseWriter := &responseWriterWrapper{
			ResponseWriter: w,
		}

		handler.ServeHTTP(
			wrappedResponseWriter,
			r,
		)

		statusCode := wrappedResponseWriter.statusCode

		data := map[string]interface{}{
			"message":     "Request finished",
			"method":      r.Method,
			"url":         r.URL.String(),
			"user_agent":  r.Header.Get("User-Agent"),
			"status_code": statusCode,
			"size":        wrappedResponseWriter.size,
			"latency":     time.Since(startTime).Milliseconds(),
		}

		if statusCode >= 400 && statusCode < 500 {
			data["message"] = "Request failed with 4xx"
			log.Warn(data)
		} else if statusCode > 500 {
			data["message"] = "Request failed with 5xx"

			log.Error(data)
		} else {
			log.Debug(data)
		}
	})
}
