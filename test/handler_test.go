package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	svcUserV1 "svc-with-grpc-gateway/api/v1"
	"svc-with-grpc-gateway/internal/handler"
	"svc-with-grpc-gateway/internal/repository"
	"svc-with-grpc-gateway/internal/service"
	"svc-with-grpc-gateway/internal/store"
)

//integration test : http1
//TODO add grpc unit test / client test

func setupHttpServerTest() http.Handler {
	db := store.NewPostgreeSQLDbTest().Connect()
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	mux := runtime.NewServeMux()
	svcUserV1.RegisterUserServiceHandlerServer(context.Background(), mux, userHandler)
	return mux
}

func TestHTTP1Create(t *testing.T) {
	r := setupHttpServerTest()
	w := httptest.NewRecorder()

	reqPayload := strings.NewReader("{\"name\":\"John Due Test\"}")
	req, _ := http.NewRequest(http.MethodPost, "/user.v1.UserService/Create", reqPayload)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	var response svcUserV1.CreateResponse
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", response.Message)
}
