package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lahaehae/crud_project/internal/db"
	"github.com/lahaehae/crud_project/internal/handler"
	"github.com/lahaehae/crud_project/internal/repository"
	"github.com/lahaehae/crud_project/internal/service"
	"github.com/lahaehae/crud_project/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	log.Printf("Starting REST server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		attribute.String("service.name", "rest-service"),
		attribute.String("service.version", "1.0.0"),
		attribute.String("service.instance.id", "instance-123"),
	)
	

	// Подключение к OpenTelemetry Collector
	// otelConn, err := telemetry.InitConn()
	// if err != nil {
	// 	log.Fatalf("Ошибка подключения к OTEL Collector: %v", err)
	// }
	// log.Println("Успешное подключение к OTEL Collector")
	// defer otelConn.Close()

	shutdownTracer, err := telemetry.InitTracerProvider(ctx, res)
	if err != nil {
		log.Fatalf("Ошибка инициализации трассировки: %v", err)
	}
	log.Println("Успешное инициализация трассировки")
	defer shutdownTracer(ctx)

	shutdownMeter, err := telemetry.InitMeterProvider(ctx, res)
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
	//dependency injection
	userRepository := repository.NewUserRepository(conn)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	

	r := gin.Default()

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)
	r.POST("/transfer", userHandler.TransferFunds)

	log.Println("Server is running on :8080")
	http.ListenAndServe("0.0.0.0:8080", r)

}