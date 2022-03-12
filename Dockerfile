FROM golang:1.17.6-alpine

WORKDIR /forum

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN apk add build-base
RUN go build -o forum 

EXPOSE 8080

CMD ["./forum"]