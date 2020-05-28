FROM golang:1.14.3

WORKDIR /go/src/github.com/fabiodmferreira/crypto-trading/coin-historical-collector
COPY . .

RUN mkdir -p /go/src/github.com/pilu/fresh
RUN go get github.com/pilu/fresh
RUN go get ./...

RUN cd cmd/serviced

CMD [ "fresh" ]
