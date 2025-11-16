FROM golang:1.24.5

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/app
RUN go build -o main main.go


EXPOSE 8080
CMD ["./main"]