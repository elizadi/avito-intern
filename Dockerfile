FROM golang:alpine
WORKDIR /
COPY . .
RUN go build -o app main.go
EXPOSE 8080
CMD ["./app"]