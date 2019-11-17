FROM umputun/baseimage:buildgo-latest as builder
#FROM golang:1.12 as builder
WORKDIR /go/simple-peak-flowmeter
COPY go.mod go.sum ./
RUN go mod download

COPY . /go/simple-peak-flowmeter
RUN go build ./cmd/web

FROM umputun/baseimage:app-latest
#FROM amd64/alpine:latest

ARG AUTHORIZED_IP_BUILDTIME
ENV AUTHORIZED_IP ${AUTHORIZED_IP_BUILDTIME}

WORKDIR /srv
COPY --from=builder /go/simple-peak-flowmeter/web /srv/web
RUN \
    chown -R app:app /srv && \
    chmod +x /srv/web
CMD /srv/web -authorizedIp=${AUTHORIZED_IP} -dsn=mongodb://mongo:27017

EXPOSE 3333