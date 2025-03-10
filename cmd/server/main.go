package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/lahaehae/crud_project/internal/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	//"go.opentelemetry.io/otel/trace"
	"github.com/lahaehae/crud_project/internal/db"
	repository "github.com/lahaehae/crud_project/internal/repository"
	service "github.com/lahaehae/crud_project/internal/service"
	"github.com/lahaehae/crud_project/internal/telemetry"
)

func main() {
	log.Printf("Waiting for connection...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		attribute.String("service.name", "grpc-service"),
		attribute.String("service.version", "1.0.0"),
		attribute.String("service.instance.id", "instance-123"),
	)
	

	// Подключение к OpenTelemetry Collector
	otelConn, err := telemetry.InitConn()
	if err != nil {
		log.Fatalf("Ошибка подключения к OTEL Collector: %v", err)
	}
	log.Println("Успешное подключение к OTEL Collector")
	defer otelConn.Close()

	shutdownTracer, err := telemetry.InitTracerProvider(ctx, res, otelConn)
	if err != nil {
		log.Fatalf("Ошибка инициализации трассировки: %v", err)
	}
	log.Println("Успешное инициализация трассировки")
	defer shutdownTracer(ctx)

	shutdownMeter, err := telemetry.InitMeterProvider(ctx, res, otelConn)
	if err != nil {
		log.Fatalf("Ошибка инициализации метрик: %v", err)
	}
	log.Println("Успешное инициализация метрик")
	defer shutdownMeter(ctx)


	telemetry.InitMetrics()
	//"postgres://postgres:postgres@localhost:5433/crud_project?sslmode=disable"
	connStr := os.Getenv("DATABASE_URL")
	//dsn := "postgres://postgres:postgres@localhost:5433/crud_project?sslmode=disable"
	conn, err := db.InitDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database")
	}

	userRepository := repository.NewUserRepository(conn)
	userService := service.NewUserService(*userRepository)

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to connect to tcp server at 9001: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.MaxConcurrentStreams(1000),	
	)

	reflection.Register(grpcServer)


	pb.RegisterUserServiceServer(grpcServer, userService)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	

}