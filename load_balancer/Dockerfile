FROM golang:alpine
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o lb ./cmd/lb
CMD [ "./lb" ]