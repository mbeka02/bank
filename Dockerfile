#Build stage
FROM golang:1.22.2-alpine3.18 AS builder
#set work dir
WORKDIR /app
COPY . .
#build application
RUN go build -o main main.go
#run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
#expose
EXPOSE 8080
#run command
CMD ["/app/main"]

