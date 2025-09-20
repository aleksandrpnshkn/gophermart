package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	logWriter struct {
		w            http.ResponseWriter
		responseData *responseData
	}
)

func (l *logWriter) Header() http.Header {
	return l.w.Header()
}

func (l *logWriter) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.responseData.size += size
	return size, err
}

func (l *logWriter) WriteHeader(statusCode int) {
	l.w.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

func newLogWriter(w *http.ResponseWriter) *logWriter {
	responseData := responseData{}

	writer := logWriter{
		w:            *w,
		responseData: &responseData,
	}

	return &writer
}

func NewLogMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()

			logWriter := newLogWriter(&res)

			next.ServeHTTP(logWriter, req)

			statusCode := logWriter.responseData.status
			responseSize := logWriter.responseData.size
			duration := time.Since(start)

			logger.Info(req.Method+" "+req.RequestURI,
				zap.Int64("duration", duration.Microseconds()),
				zap.Int("status", statusCode),
				zap.Int("size", responseSize),
			)
		})
	}
}
