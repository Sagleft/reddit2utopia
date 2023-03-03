FROM golang:1.18.3-buster as builder
WORKDIR /app

COPY ./ /go/src
WORKDIR /go/src/cmd/service
RUN go get -v -t -d ./
RUN go build -o /app/bot

FROM gcr.io/distroless/base
COPY --from=builder /app /app

CMD ["./bot"]
