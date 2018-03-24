package main

import (
	"github.com/mong0520/linebot-ptt-beauty/models"
	"gopkg.in/mgo.v2"
	"log"
	"github.com/mong0520/linebot-ptt-beauty/bots"
	"github.com/mong0520/linebot-ptt-beauty/utils"
	"os"
	"path"
)

var logger *log.Logger
var meta = &models.Model{}
var LogRoot = "logs"

func initLineBot() {
	bots.InitLineBot(meta)
}

func initDB() {
	if session, err := mgo.Dial("localhost:27017"); err != nil {
		logger.Fatalln("Unable to connect DB", err)
	} else {
		meta.Session = session
		meta.Collection = session.DB("ptt").C("beauty")
		meta.CollectionUserFavorite = session.DB("ptt").C("users")
	}
}

func main() {
	logFile, err := initLogFile()
	defer logFile.Close()

	if err != nil {
		logger.Fatalln("open file error !")
	}
	logger = utils.GetLogger(logFile)
	meta.Log = logger
	meta.Log.Println("Start to init DB...")
	initDB()
	meta.Log.Println("...Done")
	//results, _ := controllers.GetMostLike(meta.Collection, 5, 0)
	//for _, r := range results {
	//	fmt.Println(r.MessageCount.All, r.MessageCount.Boo, r.Date, r.URL, r.ArticleTitle)
	//}
	meta.Log.Println("Start to init Line Bot...")
	initLineBot()
	meta.Log.Println("...Exit")
}

func initLogFile() (logFile *os.File, err error) {
	logfilename := "pttbeauty.log"
	logFileName := path.Base(logfilename)
	logFilePath := path.Join(LogRoot, logFileName)
	if _, err := os.Stat(LogRoot); os.IsNotExist(err) {
		os.Mkdir(LogRoot, 0755)
	}
	return os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

//
//func start(){
//    initLinebot()
//    initDB()
//    results, _ := controllers.GetMostLike(meta.Collection, 5)
//    for _, r := range results{
//        fmt.Println(r.MessageCount.All, r.MessageCount.Boo, r.Date, r.URL, r.ArticleTitle)
//    }
//}
