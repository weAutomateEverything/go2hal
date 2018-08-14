FROM ubuntu:14.04

RUN apt-get update && apt-get install -y --no-install-recommends \
        g++ \
        gcc \
        libc6-dev \
        make \
        pkg-config \
        openssh-client

WORKDIR /app
# Now just add the binary
COPY go2hal /app/
COPY swagger.json /app/
COPY ./go/src/appdynamics/lib/  /app/
ENV LD_LIBRARY_PATH /app/


ENTRYPOINT ["/app/go2hal"]
EXPOSE 8000 8080 6060