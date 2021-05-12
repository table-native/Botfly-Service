FROM ubuntu:latest

# web port
EXPOSE 8080
# grpc port
EXPOSE 50051

ADD build/Botfly-Service /app/Botfly-Service
RUN ls -l

CMD /app/Botfly-Service
