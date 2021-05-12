FROM alpine:latest

RUN apk add -U --no-cache ca-certificates

WORKDIR /app
ADD build/Botfly-Service .
RUN ls -l

# web port
EXPOSE 8080
# grpc port
EXPOSE 50051
CMD ./Botfly-Service