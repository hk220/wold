FROM golang
ADD . /go/src/github.com/hk220/wold
RUN go install github.com/hk220/wold
EXPOSE 3000
ENTRYPOINT ["wold"]