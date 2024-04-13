#Build stage
FROM golang:1.22.2-alpine3.18 AS builder
#set working directory
WORKDIR /app
#copy from current dir to work dir
COPY . .
#build application
RUN go build -o main main.go
RUN apk add curl
#install goose
RUN curl -fsSL \https://raw.githubusercontent.com/pressly/goose/master/install.sh |\sh 
#run stage
FROM alpine:3.18
WORKDIR /app
#copy executables , shell files and schema to image
COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/goose ./goose
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY /sql/schema ./schema
#expose
EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
