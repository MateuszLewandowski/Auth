FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth ./cmd/main.go

FROM gcr.io/distroless/static-debian11

WORKDIR /app
COPY --from=builder ./app/auth ./app/auth
COPY --from=builder /app/.env.prod /app/.env.prod

ENTRYPOINT ["./app/auth"]
