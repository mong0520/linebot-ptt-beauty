package utils

import (
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// GetLogger: Get logger
func GetLogger(f *os.File) (logger *log.Logger) {
	if f != nil {
		logger = log.New(io.MultiWriter(os.Stdout, f), "", log.LstdFlags|log.Lshortfile)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}

	return logger
}

// CheckTitleWithBeauty: Check title with beauty
func CheckTitleWithBeauty(title string) bool {
	d, _ := regexp.MatchString("^\\[正妹\\].*", title)
	return d
}

// GetPttIDFromURL: Get ptt id from url
func GetPttIDFromURL(title string) string {
	ret := strings.TrimPrefix(title, "https://www.ptt.cc/bbs/Beauty/")
	return strings.TrimRight(ret, ".html")
}

// GetRandomIntSet: Get random int set
func GetRandomIntSet(max int, count int) (randInts []int) {
	rand.Seed(time.Now().UnixNano())
	list := rand.Perm(max)
	randInts = list[:count]
	return randInts
}

// InArray: Check if string item is in array
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

// RemoveStringItem: Remove string item from slice
func RemoveStringItem(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
