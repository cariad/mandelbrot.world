FROM golang:1.24.1 AS builder

ENV CGO_ENABLED=0

WORKDIR /build

COPY . .

RUN go mod download && \
    go build -o server .

FROM alpine

RUN apk --no-cache add curl

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /build/frontend/ /app/frontend/
COPY --from=builder /build/server    /app

CMD ["./server"]
