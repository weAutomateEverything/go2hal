FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY cacert.pem /etc/ssl/certs/ca-bundle.crt
COPY go2hal /app/
RUN apk add --no-cache openssh-client
ENTRYPOINT ["/app/go2hal"]
EXPOSE 8000 8080