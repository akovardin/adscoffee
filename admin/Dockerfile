# Build
FROM mirror.gcr.io/library/golang:1.26.0 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./

ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -v -o /admin ./cmd/admin

# Tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy
FROM mirror.gcr.io/library/debian:11-slim AS build-release-stage
ENV TZ="Europe/Moscow"
WORKDIR /

COPY --from=build-stage /admin /admin
COPY configs/* /configs/
RUN apt-get update
RUN apt-get install -y ca-certificates

EXPOSE 8080

ENTRYPOINT [ "/admin"]
