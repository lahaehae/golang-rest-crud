build_proto:
	protoc --proto_path=proto \
       --go_out=internal/pb --go_opt=paths=source_relative \
       --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
       proto/user.proto

build_grpc_ui:
	grpcui -plaintext localhost:9001


docker_otel_collector:
	docker run --rm -p 4317:4317 -p 4318:4318 \
        -p 55680:55680 -p 9464:9464 \
        -v C:/Users/aidyn/Desktop/crud_project/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
        otel/opentelemetry-collector-contrib:latest --config /etc/otel-collector-config.yaml

docker_jaeger:
	docker run --rm -d --name jaeger \
  -p 16686:16686 -p 14250:14250 -p 14317:14317 \
  jaegertracing/all-in-one:latest



docker_grafana:
	docker run -d --name=grafana -p 3000:3000 grafana/grafana

docker_prometheus:
	docker run --rm -d --name=prometheus -p 9090:9090 \
    -v C:/Users/aidyn/Desktop/crud_project/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus



