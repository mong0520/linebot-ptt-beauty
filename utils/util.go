package utils

import (
    "os"
    "log"
    "io"
    "time"
    "math/rand"
    "reflect"
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


func InArray(val interface{}, array interface{}) (exists bool, index int) {
    exists = false
    index = -1

    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                index = i
                exists = true
                return
            }
        }
    }
    return
}

func RemoveStringItem(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}
