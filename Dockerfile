FROM golang:1.21


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main cmd/main.go

EXPOSE 8050
CMD ["./main"]