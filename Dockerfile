FROM alpine:latest

RUN apk add -U --no-cache ca-certificates
ADD build/Botfly-Service .
ADD grpcwebproxy .

CMD ./Botfly-Service; ./grpcwebproxy --run_tls_server=false --allow_all_origins --backend_addr=localhost:50051