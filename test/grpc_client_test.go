package test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	svcUserV1 "svc-with-grpc-gateway/api/v1"
	"svc-with-grpc-gateway/internal/handler"
	"svc-with-grpc-gateway/internal/repository"
	"svc-with-grpc-gateway/internal/service"
	"svc-with-grpc-gateway/internal/store"
)

func setupGrpcServerTest() *bufconn.Listener {
	listener := bufconn.Listen(1024 * 1024)

	db := store.NewPostgreeSQLDbTest().Connect()
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	grpcServer := grpc.NewServer()
	svcUserV1.RegisterUserServiceServer(grpcServer, userHandler)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal().Msgf("Test server exited with error: %v", err)
		}
	}()

	return listener
}

func TestGrpcClientTest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server := setupGrpcServerTest()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return server.Dial()
	}), grpc.WithInsecure())
	defer conn.Close()

	assert.NoError(t, err)

	client := svcUserV1.NewUserServiceClient(conn)
	resp, err := client.Create(ctx, &svcUserV1.CreateRequest{Name: "John"})
	assert.NoError(t, err)
	assert.Equal(t, "OK", resp.Message)
}
