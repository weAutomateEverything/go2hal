FROM golang:1.8
WORKDIR /go/src/app
RUN git clone https://github.com/zamedic/go2hal.git .
RUN go-wrapper download
RUN go-wrapper install
ENTRYPOINT ["/go/bin/app"]
EXPOSE 8000
