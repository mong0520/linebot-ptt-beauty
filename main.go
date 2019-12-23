package main

import (
	"log"
	"os"
	"path"

	"github.com/mong0520/linebot-ptt-beauty/bots"
	"github.com/mong0520/linebot-ptt-beauty/models"
	"github.com/mong0520/linebot-ptt-beauty/utils"
	"gopkg.in/mgo.v2"
)

var logger *log.Logger
var meta = &models.Model{}
var logRoot = "logs"

func initDB(dbHostPort string) {
	if session, err := mgo.Dial(dbHostPort); err != nil {
		logger.Fatalln("Unable to connect DB", err)
	} else {
		meta.Session = session
		meta.Collection = session.DB("ptt").C("beauty")
		meta.CollectionUserFavorite = session.DB("ptt").C("users")
	}
}

func main() {
	logFile, err := initLogFile()
	dbHostPort := os.Getenv("MongoDBHostPort")
	defer logFile.Close()

	if err != nil {
		logger.Fatalln("open file error !")
	}
	logger = utils.GetLogger(logFile)
	meta.Log = logger
	meta.Log.Println("Start to init DB...", dbHostPort)
	initDB(dbHostPort)
	meta.Log.Println("...Done")

	meta.Log.Println("Start to init Line Bot...")
	bots.InitLineBot(meta, bots.ModeHTTP, "", "")
	meta.Log.Println("...Exit")
}

func initLogFile() (logFile *os.File, err error) {
	logfilename := "pttbeauty.log"
	logFileName := path.Base(logfilename)
	logFilePath := path.Join(logRoot, logFileName)
	if _, err := os.Stat(logRoot); os.IsNotExist(err) {
		os.Mkdir(logRoot, 0755)
	}
	return os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}
