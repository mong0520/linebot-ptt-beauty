package main

import (
    "gopkg.in/mgo.v2"
    "github.com/mong0520/linebot-ptt/models"
    "log"
    //"github.com/mong0520/linebot-ptt/controllers"
    //"os"
    "github.com/mong0520/linebot-ptt/bot"
)

var logger *log.Logger
var meta = &models.Model{}


func initLineBot(){
    bot.InitLineBot(meta)
}

func initDB(){
    if session, err := mgo.Dial("localhost:27017"); err != nil {
        logger.Fatalln("Unable to connect DB", err)
    } else {
        meta.Session = session
        meta.Collection = session.DB("ptt").C("beauty")
    }
}

func main() {
    initDB()
    initLineBot()
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
