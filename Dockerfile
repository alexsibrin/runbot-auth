FROM golang:1.21.6-bullseye AS build
LABEL authors="Alex Sibrin"

WORKDIR /temp

COPY . .

EXPOSE 8081/tcp \
       8088/tcp \
       5432/tcp

RUN go get ./... && go build -o bin/app ./cmd/main.go

FROM debian:12-slim
LABEL authors="Alex Sibrin"

ENV GIN_MODE=release \
    COMMON_VERSION=1.0.1 \
    COMMON_HEALTH="I'm okay" \
    RESTSERVER_HOST=0.0.0.0 \
    RESTSERVER_PORT=8081 \
    GRPCSERVER_PORT=8088 \
    POSTGRESQL_DB=runbotdb \
    POSTGRESQL_HOST=host.docker.internal \
    POSTGRESQL_PORT=5432 \
    POSTGRESQL_USER=runbot_user \
    POSTGRESQL_PASSWORD=runbotpswd \
    POSTGRESQL_SSLMODE=disable \
    POSTGRESQL_MAXOPENCONNECTIONS=4 \
    POSTGRESQL_MAXIDLECONNECTIONS=8 \
    POSTGRESQL_CONNECTIONMAXLIFETIME=200ms \
    POSTGRESQL_CONNECTIONMAXIDLETIME=5s \
    JWT_SALT=salt \
    JWT_ISSUER=iam \
    JWT_SUBJECT=auth \
    JWT_AUDIENCE="runbot users" \
    JWT_EXPIRESIN=5m \
    LOGGER_LEVEL=6 \
    LOGGER_COLORS=true \
    LOGGER_FULLTIMESTAMP=true

WORKDIR /appbin

COPY --from=build /temp/bin/app ./

CMD ["./app"]