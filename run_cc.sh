#!/bin/bash
ACTION=$1
cd /home/mong/go/src/github.com/mong0520/ChainChronicleGo
/usr/lib/go-1.10/bin/go run main.go -c conf/mong.conf -a $1