FROM golang:1.11.1 AS Builder

RUN apt-get update \
    && apt-get install -y git

RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/app
WORKDIR /go/src/app
RUN dep ensure

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags '-w -extldflags "-static"'

# For running on the machine you build on
#RUN CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"'

FROM alpine:3.8 AS Runner

RUN apk add --update ca-certificates
COPY  --from=Builder /go/src/app/app /usr/local/bin/app

CMD [ "app" ]