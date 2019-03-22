FROM golang:1.11

WORKDIR /go/src/github.com/mong0520/linebot-ptt-beauty
ADD . /go/src/github.com/mong0520/linebot-ptt-beauty

#RUN dep ensure -v
RUN GO111MODULE=on go build

ENTRYPOINT [ "./linebot-ptt-beauty" ]