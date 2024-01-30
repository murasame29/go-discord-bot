FROM golang:1.21.4 AS builder
WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM ubuntu:18.04 AS runner
WORKDIR /app

COPY --from=builder /build/main .

ENV DISCORD_TOKEN=

CMD ["./main"]