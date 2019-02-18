FROM golang:1.11

WORKDIR /twitchstream
COPY . .

RUN go get -d -v ./...
RUN go build github.com/seregayoga/twitchstream/cmd/twitchstream

CMD ["/twitchstream/twitchstream"]