FROM golang:1.18

ADD ./ /API

WORKDIR /API
RUN go mod init Hotel
RUN go mod tidy
RUN go build -o Hotel .
RUN go test ./...
RUN chmod +x /API

ENTRYPOINT ["/API/Hotel"]
