FROM golang:alpine

RUN apk update && apk add --no-cache git

ARG GROUP_ID
ARG USER_ID

ENV GIN_PORT 8080
ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN go get github.com/codegangsta/gin

WORKDIR $GOPATH/src/dns/
RUN echo $GOPATH

EXPOSE $GIN_PORT
CMD ["gin", "-t", "./cmd/api", "run"]