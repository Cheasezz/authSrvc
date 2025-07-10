FROM golang:1.24.4-alpine3.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o authSrvc cmd/mainapp/mainapp.go
RUN CGO_ENABLED=0 GOOS=linux go build -o authSrvcMigrate cmd/migrator/migrator.go

FROM alpine:latest AS build-release

WORKDIR /

COPY --from=build /app/authSrvc .
COPY --from=build /app/authSrvcMigrate .
COPY --from=build /app/config/config.yml config/
COPY --from=build /app/migrations migrations/

EXPOSE 8080

CMD ["/authSrvc"]