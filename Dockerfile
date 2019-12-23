FROM golang:latest

WORKDIR $GOPATH/src/mong0520/linebot-ptt-beauty
COPY . $GOPATH/src/mong0520/linebot-ptt-beauty
RUN GO111MODULE=on go build

EXPOSE 5000
ENTRYPOINT ["./linebot-ptt-beauty"]
CMD ["./linebot-ptt-beauty"]