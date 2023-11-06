FROM golang:1.18

ADD ./ /API

WORKDIR /API
RUN go mod init Search
RUN go mod tidy
RUN go build -o Search .
RUN chmod +x /API

ENTRYPOINT ["/API/Search"]
