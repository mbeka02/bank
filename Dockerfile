#Build stage
FROM golang:1.22.2-alpine3.18 AS builder
#set working directory
WORKDIR /app
#copy from current dir to work dir
COPY . .
#build application
RUN go build -o main main.go
#run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
#expose
EXPOSE 8080
#run command
CMD ["/app/main"]

