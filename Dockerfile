FROM alpine:latest

RUN apk add -U --no-cache ca-certificates

WORKDIR /app
ADD build/Botfly-Service .
RUN ls -l

CMD ./Botfly-Service