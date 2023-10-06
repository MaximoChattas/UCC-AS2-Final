FROM golang:1.18

ADD ./ /API

WORKDIR /API
RUN go mod init User-Reservation
RUN go mod tidy
RUN go build -o User-Reservation .
RUN go test ./...
RUN chmod +x /API

ENTRYPOINT ["/API/User-Reservation"]
