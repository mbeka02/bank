FROM golang:1.22.2-alpine3.18
#set work dir
WORKDIR /app
COPY . .
#build application
RUN go build -o bank bank.go
#expose
EXPOSE 8080
#run command
CMD ["/app/bank"]

