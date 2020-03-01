FROM golang as builder

WORKDIR /build

COPY . .

RUN GCO_ENABLED=0 GOOS=linux go build -o gcp-ddns-updater


FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /build/gcp-ddns-updater /gcp-ddns-updater
RUN chmod +x /gcp-ddns-updater

CMD ["/gcp-ddns-updater"]
