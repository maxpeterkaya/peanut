FROM debian:bookworm AS certs

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

FROM scratch

# set env to tell the binary to use absolute paths
ENV IS_CONTAINER=true

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY peanut /usr/bin/peanut

ENTRYPOINT ["/usr/bin/peanut"]