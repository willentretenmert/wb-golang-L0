FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-server .

FROM alpine

RUN apk add --no-cache bash

COPY --from=builder /app/go-server /go-server
COPY --from=builder /app/template.html /template.html
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

ENTRYPOINT ["/wait-for-it.sh", "postgres:5432", "--", "/wait-for-it.sh", "stan:4222", "--"]

CMD ["/go-server"]
