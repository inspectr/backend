FROM golang:alpine

ARG APP_NAME=inspectr
ARG APP_PATH=/go/src/github.com/inspectr/backend

RUN apk -U add git
RUN go get github.com/cespare/reflex
RUN go get -u github.com/jteeuwen/go-bindata/...

RUN mkdir -p $APP_PATH
WORKDIR $APP_PATH

COPY . $APP_PATH

RUN go-bindata -pkg assets -o assets/assets.go \
		plugins/api/schema.graphql \
		plugins/api/static/

RUN /go/bin/go-bindata -pkg assets -o assets/assets.go plugins/api/schema.graphql
RUN go build -i -o /go/bin/$APP_NAME .
