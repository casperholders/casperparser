# Docker builder for Golang
FROM golang as builder

WORKDIR /go/src/github.com/casperholders/casperparser
COPY . .
RUN set -x && \
    touch .env && \
    make build
# Docker run Golang app
FROM scratch

WORKDIR /root/
COPY --from=builder /go/src/github.com/casperholders/casperparser/bin/casperParser .
COPY .blank .casperParser.yaml? /root/
ENTRYPOINT ["./casperParser"]