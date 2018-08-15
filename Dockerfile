# Container to build the Go binary
FROM golang as builder
WORKDIR /go/src/github.com/tinyzimmer/s3explorer/
ADD . .
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o app


# Container which contains binary
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/tinyzimmer/s3explorer/app /bin/s3explorer
ENTRYPOINT ["/bin/s3explorer"]


# To Run:
# docker build -t s3explorer .
# docker run -it --rm -v ~/.aws/credentials:/root/.aws/credentials:ro s3explorer
