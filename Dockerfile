FROM golang:1.24.1

RUN go version

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x ./wait-for-postgres.sh

RUN go mod download
RUN go build -o music-library-server-app ./cmd/music-library/main.go
CMD ["./music-library-server-app"]