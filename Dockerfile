FROM alpine:3.6
WORKDIR /app
# Now just add the binary
#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY cacert.pem /etc/ssl/certs/ca-bundle.crt
COPY go2hal /app/
ENTRYPOINT ["/app/go2hal"]
EXPOSE 8000