FROM golang:1.17-alpine AS builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .

FROM alpine:3 as certs
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=builder /main ./
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["./main"]
EXPOSE 8080