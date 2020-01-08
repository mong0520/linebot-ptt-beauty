package main

import (
	"crypto/tls"
	"log"
	"net"
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

func initDB(dbURI string) {
	dialInfo, _ := mgo.ParseURL(dbURI)
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	if session, err := mgo.DialWithInfo(dialInfo); err != nil {
		logger.Fatalln("Unable to connect DB", err)
	} else {
		meta.Session = session
		meta.Collection = session.DB("ptt").C("beauty")
		meta.CollectionUserFavorite = session.DB("ptt").C("users")
	}
}

func main() {
	logFile, err := initLogFile()
	// dbHostPort := os.Getenv("MongoDBHostPort")
	dbURI := os.Getenv("MongoDBURI")
	defer logFile.Close()

	if err != nil {
		logger.Fatalln("open file error !")
	}
	logger = utils.GetLogger(logFile)
	meta.Log = logger
	meta.Log.Println("Start to init DB...", dbURI)
	initDB(dbURI)
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
