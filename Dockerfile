FROM golang

RUN go get github.com/rudes/runestats
RUN go install github.com/rudes/runestats
ENTRYPOINT /go/bin/runestats

EXPOSE 8080