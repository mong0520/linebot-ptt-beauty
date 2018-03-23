package utils

import (
    "os"
    "log"
    "io"
    "time"
    "math/rand"
)

func GetLogger(f *os.File)(logger *log.Logger){
    if f != nil{
        logger = log.New(io.MultiWriter(os.Stdout, f), "", log.LstdFlags | log.Lshortfile)
    }else{
        logger = log.New(os.Stdout, "", log.LstdFlags | log.Lshortfile)
    }

    return logger
}

func GetRandomIntSet(max int, count int)(randInts []int){
    rand.Seed(time.Now().UnixNano())
    list := rand.Perm(max)
    randInts = list[:count]
    return randInts
}