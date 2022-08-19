package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	svcUserV1 "svc-with-grpc-gateway/api/v1"
	"svc-with-grpc-gateway/config"
	"svc-with-grpc-gateway/internal/handler"
	"svc-with-grpc-gateway/internal/middleware"
	"svc-with-grpc-gateway/internal/repository"
	"svc-with-grpc-gateway/internal/router"
	"svc-with-grpc-gateway/internal/service"
	"svc-with-grpc-gateway/internal/store"
	"svc-with-grpc-gateway/pkg/observability"
)

type server struct {
	mux        http.Handler
	grpcServer *grpc.Server

	stopPusherMetricsFn      observability.StopPusherFunc
	shutdownTracerExporterFn observability.ShutDownFunc
}

func NewServer() *server {
	//metrics & tracer init
	stopPusher := observability.InitMetricProvider(config.OtlpCollectorUrl)
	shutDownTracer := observability.InitTracerProvider(config.OtlpCollectorUrl, config.ServiceName, true)

	db := store.NewSQLLite()
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	//register grpc
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	svcUserV1.RegisterUserServiceServer(grpcServer, userHandler)

	//register grpc-health
	grpcHealthHandler := handler.NewHealthCheckServerHandler()
	grpc_health_v1.RegisterHealthServer(grpcServer, grpcHealthHandler)

	//register grpc-gateway
	gwMux := runtime.NewServeMux()
	err := svcUserV1.RegisterUserServiceHandlerServer(context.Background(), gwMux, userHandler)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return nil
	}

	mux := router.NewRouter()
	//adding custom route http1
	mux.Handle("/", gwMux)

	return &server{
		mux:                      mux,
		grpcServer:               grpcServer,
		stopPusherMetricsFn:      stopPusher,
		shutdownTracerExporterFn: shutDownTracer,
	}
}

func (s *server) Run() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	serv := http.Server{
		Addr:    config.Address,
		Handler: middleware.GrpcHttpMiddleware(s.grpcServer, otelhttp.NewHandler(s.mux, "svc-with-grpc-gateway")),
	}

	go func() {
		err := serv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msg(err.Error())
		}
	}()

	//wait signal interrupt
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := serv.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return
	}

	err = s.stopPusherMetricsFn(ctx)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return
	}

	err = s.shutdownTracerExporterFn(ctx)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return
	}

	log.Info().Msg("server exited properly")
}
