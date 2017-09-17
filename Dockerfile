FROM golang:alpine
WORKDIR /go/src/app
COPY . .
RUN go-wrapper download
RUN go-wrapper install
ENTRYPOINT ["/go/bin/app"]
EXPOSE 8000
