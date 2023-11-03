FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-server .

FROM alpine

COPY --from=builder /app/go-server /go-server

ENTRYPOINT ["/go-server"]