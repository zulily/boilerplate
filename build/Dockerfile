FROM golang:1.4.2-cross
MAINTAINER Core Engineering <core@zulily.com>

RUN go get github.com/tools/godep
RUN go get -u github.com/golang/lint/golint

# We're building static binaries here!
ENV CGO_ENABLED 0

# A volume mount can be made to /output in order to retrieve the compiled binary
ENV GOBIN /output

ENV GIT_SHA unknown

# Invokes a shell explicitly, to get param expansion of $GIT_SHA and $USER
CMD ["sh", "-c", "godep go build -o ${BINARY} -installsuffix cgo -a -tags netgo -ldflags -s -ldflags \"-X main.BuildSHA ${GIT_SHA}\" "]
