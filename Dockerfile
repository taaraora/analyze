FROM golang:latest as builder

WORKDIR $GOPATH/src/github.com/supergiant/robot/cmd/analyzed
COPY . $GOPATH/src/github.com/supergiant/robot/

ARG ARCH=amd64
ARG GO111MODULE=on

RUN go get -v -d ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s' -o /go/bin/analyzed

FROM scratch
COPY --from=builder /go/bin/analyzed /bin/analyzed

ENTRYPOINT ["/bin/analyzed"]
