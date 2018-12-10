FROM golang:1-alpine
RUN apk --no-cache add ca-certificates

ENV APP_PATH=github.com/mat285/advent-bot/
ENV APP_ROOT=/go/src/${APP_PATH}

ADD . ${APP_ROOT}/.

RUN go install ${APP_PATH}
ENTRYPOINT ["/go/bin/advent-bot"]
