package router

import "net/http"

func NewRouter() *http.ServeMux {
	//TODO add tracer and metrics middleware
	//TODO add healthz route

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})

	return mux
}
