FROM golang:1.13-alpine

COPY . /work

WORKDIR /work

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prometheus-nginxlog-exporter

FROM alpine

RUN apk add --no-cache -X http://dl-cdn.alpinelinux.org/alpine/edge/testing watchexec

COPY --from=0 /work/prometheus-nginxlog-exporter /prometheus-nginxlog-exporter
ADD entrypoint.sh /entrypoint.sh

EXPOSE 4040

ENTRYPOINT ["/entrypoint.sh"]