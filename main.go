package main

import (
	"log"
	"os"
	"path"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/kkdai/linebot-ptt-beauty/bots"
	"github.com/kkdai/linebot-ptt-beauty/controllers"
	"github.com/kkdai/linebot-ptt-beauty/models"
	"github.com/kkdai/linebot-ptt-beauty/utils"
)

var logger *log.Logger
var meta = &models.Model{}
var logRoot = "logs"

func main() {
	logFile, err := initLogFile()
	// dbHostPort := os.Getenv("MongoDBHostPort")
	defer logFile.Close()

	options, _ := pg.ParseURL(os.Getenv("DATABASE_URL"))
	db := pg.Connect(options)
	meta.Db = db
	defer db.Close()

	err = createSchema(db)
	if err != nil {
		panic(err)
	}

	users := []controllers.UserFavorite{}
	err = db.Model(&users).Select()
	if err != nil {
		log.Println(err)
	}
	log.Println("***Start server all users =", users)
	if err != nil {
		logger.Fatalln("open file error !")
	}
	logger = utils.GetLogger(logFile)
	meta.Log = logger

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

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*controllers.UserFavorite)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
