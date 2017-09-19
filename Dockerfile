FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY go2hal /app/
ENTRYPOINT ["/app/go2hal"]
EXPOSE 8000