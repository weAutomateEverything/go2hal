FROM alpine
WORKDIR /app
# Now just add the binary
ADD /home/travis/gopath/bin/go2hal /app/
ENTRYPOINT ["./go2hal"]
