FROM golang:1.16 as BUILD

WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0
RUN go build .

FROM alpine:latest as SSL-CERTS

FROM scratch
COPY --from=SSL-CERTS /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=BUILD /build/contracts-dump /bin/contracts-dump
COPY .docker/group /etc/group
COPY .docker/passwd /etc/passwd
COPY --chown=contracts-dump:contracts-dump contracts/Source /var/lib/contracts-dump/contracts/Source
USER contracts-dump
WORKDIR /var/lib/contracts-dump
ENTRYPOINT [ "/bin/contracts-dump" ]
CMD [ ]