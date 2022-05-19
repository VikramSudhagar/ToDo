FROM golang:1.18.0-bullseye
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go mod download
EXPOSE 8081
RUN go build -o todo /app/app.go
CMD ["/app/todo"]