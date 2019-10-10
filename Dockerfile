FROM golang:alpine

ENV GO111MODULE on

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk --no-cache add gcc g++ make ca-certificates

WORKDIR $GOPATH/src/go-api
COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/go-api .

ENTRYPOINT  ["/go/go-api"]