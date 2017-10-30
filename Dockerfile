FROM golang:1.8

WORKDIR /go/src/app
COPY *.go ./

RUN go-wrapper download
RUN go-wrapper install

ENTRYPOINT ["/go/bin/app"]

CMD ["go-wrapper", "run"] 

