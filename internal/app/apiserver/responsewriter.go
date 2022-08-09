package apiserver

import "net/http"

type responseWriter struct {
	http.ResponseWriter // анонимное поле -> все его методы будут доступны внутри структуры
	code                int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
