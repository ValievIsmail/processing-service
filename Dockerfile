FROM golang:alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -mod=vendor -v -o ./proccesing-service cmd/*.go
EXPOSE 8080
CMD ["./proccesing-service"]