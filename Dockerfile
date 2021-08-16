# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.16.23 as builder

ENV GO111MODULE on
ENV GOPROXY http://proxy.golang.org

RUN mkdir -p /src/github.com/StarsiegePlayers
WORKDIR /src/github.com/StarsiegePlayers

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

ADD . .
RUN buffalo build -o /bin/api

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/api .

# Uncomment to run the binary in "production" mode:
ENV GO_ENV=production

# Bind the api to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 5050

# Uncomment to run the migrations before running the binary:
# CMD /bin/api migrate; /bin/api
CMD exec /bin/api

FROM scratch AS export-stage
COPY --from=builder /bin/api /
