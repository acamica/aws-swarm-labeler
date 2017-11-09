FROM golang:1.8-alpine as builder

WORKDIR /go/src/app
COPY *.go ./

RUN apk add -U git &&\
    go-wrapper download &&\
    go-wrapper install &&\
    go build aws_swarm_labeler.go

FROM alpine

COPY --from=builder /go/src/app/aws_swarm_labeler /go/src/app/aws_swarm_labeler

RUN chmod +x /go/src/app/aws_swarm_labeler

ENTRYPOINT ["/go/src/app/aws_swarm_labeler"]

CMD ["go-wrapper", "run"]