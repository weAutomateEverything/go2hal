FROM frolvlad/alpine-glibc:latest
WORKDIR /app
# Now just add the binary
COPY go2hal /app/
COPY swagger.json /app/
COPY ./go/src/appdynamics/lib/libappdynamics.so  /usr/glibc-compat/lib/
RUN wget "https://www.archlinux.org/packages/core/x86_64/zlib/download" -O /tmp/libz.tar.xz \
    && mkdir -p /tmp/libz \
    && tar -xf /tmp/libz.tar.xz -C /tmp/libz \
    && cp /tmp/libz/usr/lib/libz.so.1.2.11 /usr/glibc-compat/lib \
    && /usr/glibc-compat/sbin/ldconfig \
    && rm -rf /tmp/libz /tmp/libz.tar.xz
RUN apk add --no-cache openssh-client
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ENTRYPOINT ["/app/go2hal"]
EXPOSE 8000 8080 6060