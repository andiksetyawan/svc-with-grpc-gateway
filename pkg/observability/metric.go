package observability

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type StopPusherFunc func(ctx context.Context) error

func InitMetricProvider(otelCollectorURL string) StopPusherFunc {
	ctx := context.Background()
	metricClient := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otelCollectorURL))

	metricExporter, err := otlpmetric.New(ctx, metricClient)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create the collector metric exporter")
	}

	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(histogram.WithExplicitBoundaries([]float64{5, 10, 25, 50, 100, 200, 400, 800, 1000})),
			metricExporter,
		),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(10*time.Second),
	)

	if err := pusher.Start(context.Background()); err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	global.SetMeterProvider(pusher)
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := pusher.Stop(ctx); err != nil {
			//log.Debug().Err(err).Msg("failed to stop metrics pusher")
			return err
		}
		log.Debug().Msg("metrics pusher has been stopped")
		return nil
	}
}
