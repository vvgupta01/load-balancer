FROM golang:bullseye

WORKDIR loadbalancer

COPY . .
RUN go mod download
RUN go install src/main.go

ENTRYPOINT main
