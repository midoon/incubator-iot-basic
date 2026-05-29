FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend-incubator ./cmd/main.go

# Stage 2: minimal runtime
FROM scratch

WORKDIR /app

COPY --from=builder /app/backend-incubator .
COPY .env .

EXPOSE 8080

CMD ["./backend-incubator"]
