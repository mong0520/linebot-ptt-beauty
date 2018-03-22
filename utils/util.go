package utils

import (
    "os"
    "log"
    "io"
)

func GetLogger(f *os.File)(logger *log.Logger){
    if f != nil{
        logger = log.New(io.MultiWriter(os.Stdout, f), "", log.LstdFlags | log.Lshortfile)
    }else{
        logger = log.New(os.Stdout, "", log.LstdFlags | log.Lshortfile)
    }

    return logger
}