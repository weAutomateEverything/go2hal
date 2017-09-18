FROM alpine:latest
WORKDIR /app
# Now just add the binary
COPY /home/travis/gopath/bin/go2hal /app/
ENTRYPOINT ["./go2hal"]
