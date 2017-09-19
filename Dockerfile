FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY go2hal /app/
RUN chmod +x /app/go2hal
RUN ls -al
ENTRYPOINT ./go2hal
EXPOSE 8000