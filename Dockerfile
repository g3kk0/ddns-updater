FROM golang as builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ddns-updater


FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /build/ddns-updater /ddns-updater
RUN chmod +x /ddns-updater

CMD ["/ddns-updater"]
