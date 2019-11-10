FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git # need git for fetching deps
COPY . $GOPATH/src/github.com/adrianosela/padl/api
WORKDIR $GOPATH/src/github.com/adrianosela/padl/api
RUN go get -d -v
RUN go build -o /go/bin/padl_api

FROM alpine:latest AS api
RUN apk add --update bash curl && rm -rf /var/cache/apk/*
COPY --from=builder /go/bin/padl_api /bin/padl_api
EXPOSE 443
EXPOSE 80
CMD ["/bin/padl_api"]
