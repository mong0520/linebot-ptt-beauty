FROM golang:1.11

WORKDIR /go/src/github.com/mong0520/linebot-ptt-beauty
ADD . /go/src/github.com/mong0520/linebot-ptt-beauty

#RUN go get -u github.com/golang/dep/...
#RUN dep ensure -v
#RUN go build

