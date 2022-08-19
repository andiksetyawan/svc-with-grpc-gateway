package main

import "svc-with-grpc-gateway/internal/server"

func main() {
	server.NewServer().Run()
}
