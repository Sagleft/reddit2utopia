FROM golang:1.18.3-buster as builder
WORKDIR /app
COPY . .
RUN go build -o /app/bot

FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /app /app

CMD ["./bot"]
