FROM alpine

#https://github.com/zenazn/goji/issues/126
RUN apk add -U ca-certificates

COPY aws_swarm_labeler /

CMD /aws_swarm_labeler

