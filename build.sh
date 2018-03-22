#!/bin/bash

GOOS=linux GOARCH=amd64 go build
scp linebot-ptt dev:/home/mong/linebot
