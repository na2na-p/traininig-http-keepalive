FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./server.go

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /server /

USER nonroot:nonroot

ENTRYPOINT ["/server"]