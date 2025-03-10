package telemetry

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

var (
	once			 sync.Once
	Meter            metric.Meter
	RequestsCounter  metric.Int64Counter
	LatencyRecorder  metric.Float64Histogram
	ErrorCounter     metric.Int64Counter
	RepoLatencyRecorder metric.Float64Histogram
)

// Initializes an OTLP exporter, and configures the corresponding meter provider.
func InitMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}

func InitMetrics() {
	once.Do(func(){
		Meter = otel.Meter("grpc-server")
		var err error
	RequestsCounter, err = Meter.Int64Counter(
		"count",
		metric.WithDescription("Общее количество gRPC-запросов"),
	)
	if err != nil {
		log.Printf("Ошибка создания счетчика запросов: %v", err)
	}

	LatencyRecorder, err = Meter.Float64Histogram(
		"latency",
		metric.WithDescription("Время обработки gRPC-запросов"),
	)
	if err != nil {
		log.Printf("Ошибка создания гистограммы задержек: %v", err)
	}

	RepoLatencyRecorder, err = Meter.Float64Histogram(
		"repository_latency",
		metric.WithDescription("Время обработки запросов в repository"),
	)
	if err != nil {
		log.Printf("failed to create repository_latency histogram")
	}

	ErrorCounter, err = Meter.Int64Counter(
		"grpc_server_errors_total",
		metric.WithDescription("Количество ошибок gRPC-запросов"),
	)
	if err != nil {
		log.Printf("Ошибка создания счетчика ошибок")
	}

	})
	
}

func RecordErrorMetric(ctx context.Context, operation string, err error) {
	if err == nil{
		return
	}

    if ErrorCounter != nil {
        ErrorCounter.Add(ctx, 1, 
            metric.WithAttributes(
                attribute.String("operation", operation),
                attribute.String("error.type", fmt.Sprintf("%T", err)),
                attribute.String("error.msg", err.Error()),
            ),
        )
    }
}