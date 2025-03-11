package telemetry

// import (
// 	"fmt"
// 	"net/http"
// 	//"google.golang.org/grpc"
// 	//"google.golang.org/grpc/credentials/insecure"
// )
// func InitConn() (*http.Client, error) {
// 	// It connects the OpenTelemetry Collector through local gRPC connection.
// 	// You may replace `localhost:4317` with your endpoint.
// 	conn, err := gg.NewClient("otel-collector:4318",
// 		// Note the use of insecure transport here. TLS is recommended in production.
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
// 	}

// 	return conn, err
// }