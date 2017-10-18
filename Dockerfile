FROM golang:1.8

WORKDIR /go/src/app
COPY ./src/ ./

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"] 

