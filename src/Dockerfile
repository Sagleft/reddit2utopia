FROM golang:1.19-alpine

COPY ./ /go/src
WORKDIR /go/src
RUN go get -v -t -d ./...
RUN go build -o /app/bot

WORKDIR /app
RUN mkdir /app/cache

CMD ["./bot"]
