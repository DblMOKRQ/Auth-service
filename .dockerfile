FROM golang:1.23.4-alpine3.21

WORKDIR /

COPY . .
RUN go mod download

RUN apk add build-base
RUN CGO_ENABLED=1  go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o /auth-service /cmd/main.go

EXPOSE 50051

CMD ["/auth-service"]