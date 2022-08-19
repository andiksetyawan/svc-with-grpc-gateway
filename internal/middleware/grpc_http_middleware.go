package middleware

import (
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

//GrpcHttpMiddleware middleware
func GrpcHttpMiddleware(nextGrpcServer *grpc.Server, nextOtherHandler http.Handler) http.Handler {
	//Package h2c implements the unencrypted "h2c" form of HTTP/2.
	//The h2c protocol is the non-TLS version of HTTP/2 which is not available from net/http or golang.org/x/net/http2.
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			nextGrpcServer.ServeHTTP(w, r)
		} else {
			nextOtherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
