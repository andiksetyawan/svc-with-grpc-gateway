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
	"svc-with-grpc-gateway/internal/router"
	"svc-with-grpc-gateway/pkg/observability"
)

type server struct {
	handler http.Handler

	stopPusherMetricsFn      observability.StopPusherFunc
	shutdownTracerExporterFn observability.ShutDownFunc
}

func NewServer() *server {
	//init config
	config.Init()
	//metrics & tracer init
	stopPusher := observability.InitMetricProvider(config.App.OtlpCollectorUrl)
	shutDownTracer := observability.InitTracerProvider(config.App.OtlpCollectorUrl, config.App.ServiceName, config.App.InsecureOtlpCollector)
	userHandlerServer := InitializedUserServiceHandlerServer()

	//register grpc
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	svcUserV1.RegisterUserServiceServer(grpcServer, userHandlerServer)

	//register grpc-health
	grpcHealthHandler := handler.NewHealthCheckServerHandler()
	grpc_health_v1.RegisterHealthServer(grpcServer, grpcHealthHandler)

	//register grpc-gateway
	gwMux := runtime.NewServeMux()
	err := svcUserV1.RegisterUserServiceHandlerServer(context.Background(), gwMux, userHandlerServer)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return nil
	}

	//adding custom route http1
	mux := router.NewRouter()
	mux.Handle("/", gwMux)

	handlerAdapter := middleware.GrpcHttpMiddleware(grpcServer, otelhttp.NewHandler(mux, "svc-with-grpc-gateway"))
	return &server{
		handler:                  handlerAdapter,
		stopPusherMetricsFn:      stopPusher,
		shutdownTracerExporterFn: shutDownTracer,
	}
}

func (s *server) Run() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	serv := http.Server{
		Addr:    config.App.Address,
		Handler: s.handler,
	}

	go func() {
		log.Debug().Msgf("listening and serving on %s", config.App.Address)
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
