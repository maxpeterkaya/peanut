FROM scratch

# set env to tell the binary to use absolute paths
ENV IS_CONTAINER=true

COPY peanut /usr/bin/peanut

ENTRYPOINT ["/usr/bin/peanut"]