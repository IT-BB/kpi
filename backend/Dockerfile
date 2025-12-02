FROM golang:1.24 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/enterprise-kpi ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /bin/enterprise-kpi /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]
