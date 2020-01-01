FROM golang:1.13-buster AS builder

ARG NETRC_ENABLED="false"
ARG NETRC_MACHINE="github.com"
ARG NETRC_LOGIN
ARG NETRC_PASSWORD

RUN test "$NETRC_ENABLED" && printf "machine ${NETRC_MACHINE}\nlogin ${NETRC_LOGIN}\npassword ${NETRC_PASSWORD}\n" >> /root/.netrc \
	&& chmod 600 /root/.netrc

RUN apt update \
	&& apt install --no-install-recommends -y \
		ca-certificates \
		libasound2-dev \
		gcc \
		git \
		pkg-config

COPY go.mod go.sum /go/src/github.com/juli3nk/doorbell-sound/
WORKDIR /go/src/github.com/juli3nk/doorbell-sound

ENV GO111MODULE on
RUN go mod download

COPY . .

RUN go install


FROM debian:stable-slim

COPY --from=builder /go/bin/doorbell-sound /usr/local/bin/doorbell-sound
COPY doorbell.mp3 /usr/local/share/doorbell/doorbell.mp3

RUN apt update \
	&& apt install --no-install-recommends -y \
		libasound2

ENTRYPOINT ["/usr/local/bin/doorbell-sound"]
