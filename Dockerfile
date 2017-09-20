FROM golang:alpine as builder
COPY . /go/src/github.com/concourse/archive-resource
ENV CGO_ENABLED 0
RUN go build -o /assets/out github.com/concourse/archive-resource/out
RUN go build -o /assets/in github.com/concourse/archive-resource/in
RUN go build -o /assets/check github.com/concourse/archive-resource/check
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM alpine:edge AS resource

RUN apk --no-cache add \
  bash \
  curl \
  gzip \
  jq \
  tar \
  openssl
COPY --from=builder /assets /opt/resource

FROM resource AS tests
COPY --from=builder /tests /tests
RUN set -e; for test in /tests/*.test; do \
		$test; \
	done

FROM resource
