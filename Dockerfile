FROM golang:1.20.4

WORKDIR /usr/src/app
COPY . .
RUN go mod tidy