# build analytics binary
FROM golang:1.24 AS build
WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -v -o /analytics ./analytics/cmd/analytics

# Deploy
FROM alpine:3.21

RUN apk add --no-cache wget

RUN wget -q https://github.com/dagu-org/dagu/releases/download/v2.16.6/dagu_2.16.6_linux_amd64.tar.gz -O /tmp/dagu.tar.gz \
  && tar -xzf /tmp/dagu.tar.gz -C /usr/local/bin/ dagu \
  && chmod +x /usr/local/bin/dagu \
  && rm /tmp/dagu.tar.gz

COPY --from=build /analytics /analytics

COPY analytics/dagu/dags/requests.yaml /var/lib/dagu/dags/requests.yaml
COPY analytics/dagu/dags/impressions.yaml /var/lib/dagu/dags/impressions.yaml
COPY analytics/dagu/dags/clicks.yaml /var/lib/dagu/dags/clicks.yaml
COPY analytics/dagu/dags/responses.yaml /var/lib/dagu/dags/responses.yaml

ENV DAGU_HOST=0.0.0.0
ENV DAGU_PORT=8080
ENV DAGU_DAG_DIR=/var/lib/dagu/dags
ENV DAGU_DATA_DIR=/var/lib/dagu/data
ENV DAGU_LOG_DIR=/var/lib/dagu/logs

RUN mkdir -p /var/lib/dagu/data /var/lib/dagu/logs

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/dagu"]
CMD ["server", "--host=0.0.0.0", "--port=8080", "--dags=/var/lib/dagu/dags"]