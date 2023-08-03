## load testing and benchmarking gRPC services
https://ghz.sh/docs/intro

ghz --proto=grpc_service.proto --call=sonyflake.SonyflakeService/NextID --insecure --async --concurrency=50 --total=10000 --data='{"Num": 1}' 0.0.0.0:8083

## build example:
go build -x -o sonyflake-grpc