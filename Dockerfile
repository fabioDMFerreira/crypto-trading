FROM golang:1.14.4

# Create app directory
WORKDIR /usr/src/app

# Bundle app source
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go get github.com/codegangsta/gin

WORKDIR /usr/src/app/cmd/webserver

CMD ["gin","-p","5000","-a","4000","--all","run","main"]
