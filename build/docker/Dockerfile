# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the actual server from the cmd/server folder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/app

# -------- STAGE 2: Run --------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]

