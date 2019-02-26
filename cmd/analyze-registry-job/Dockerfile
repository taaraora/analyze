FROM golang:1.11.5 as job_builder

ARG ARCH=amd64
ARG GO111MODULE=on

WORKDIR $GOPATH/src/github.com/supergiant/analyze/

COPY go.mod go.sum $GOPATH/src/github.com/supergiant/analyze/
RUN go mod download

COPY . $GOPATH/src/github.com/supergiant/analyze/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -o $GOPATH/bin/analyze-registry-job -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s'  ./cmd/analyze-registry-job

FROM busybox:latest
COPY --from=job_builder /go/bin/analyze-registry-job /bin/analyze-registry-job

ENTRYPOINT ["/bin/analyze-registry-job"]
