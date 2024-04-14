FROM golang:1.21 as builder

WORKDIR /app
COPY ./ ./
RUN go mod tidy
RUN go build -o /app/build/server cmd/server/main.go

FROM debian:12
COPY --from=builder /app/build/server /server
ENTRYPOINT /server