FROM golang:alpine
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o node ./cmd/node
CMD ["./node"]