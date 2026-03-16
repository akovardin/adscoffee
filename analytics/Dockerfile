# build
FROM mirror.gcr.io/library/golang:1.26.0 AS build
WORKDIR /app
COPY go.mod go.sum ./

ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -v -o /analytics ./analytics/cmd/analytics

# Deploy
FROM --platform=linux/amd64 ghcr.io/dagu-org/dagu:latest
COPY --from=build /analytics /analytics
COPY analytics/dagu/dags/requests.yaml /var/lib/dagu/dags/requests.yaml
COPY analytics/dagu/dags/impressions.yaml /var/lib/dagu/dags/impressions.yaml
COPY analytics/dagu/dags/clicks.yaml /var/lib/dagu/dags/clicks.yaml
COPY analytics/dagu/dags/responses.yaml /var/lib/dagu/dags/responses.yaml