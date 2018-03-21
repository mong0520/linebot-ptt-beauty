package controllers

import (
    "gopkg.in/mgo.v2/bson"
    "github.com/mong0520/linebot-ptt/models"
    "gopkg.in/mgo.v2"
    "fmt"
)

func GetOne(collection *mgo.Collection, query bson.M)(result *models.ArticleDocument, err error){
    //query := bson.M{"article_id": "M.1521548086.A.DCA"}
    document := &models.ArticleDocument{}
    result, err = document.GeneralQueryOne(collection, query)
    if err != nil{
        //fmt.Println(err)
        return nil, err
    }else{
        return result, nil
    }
}

func GetAll(collection *mgo.Collection, query bson.M)(results []models.ArticleDocument, err error){
    document := &models.ArticleDocument{}
    results, err = document.GeneralQueryAll(collection, query, "", -1)
    if err != nil{
        //fmt.Println(err)
        return nil, err
    }else{
        fmt.Printf("%+v", results)
        return results, nil
    }
}

func GetMostLike(collection *mgo.Collection, count int)(results []models.ArticleDocument, err error){
    document := &models.ArticleDocument{}
    //query := bson.M{"message_count.all": bson.M{"$gt": like}, "ArticleTitle": "/正妹/"}
    //query := bson.M{"ArticleTitle": bson.RegEx{"*", ""}}
    //query := bson.M{"ArticleTitle": bson.RegEx{".+", ""}}
    query := bson.M{"article_title": bson.M{"$regex": bson.RegEx{".*正妹.*", ""}}}
    results, err = document.GeneralQueryAll(collection, query, "-message_count.push", count)
    if err != nil{
        //fmt.Println(err)
        return nil, err
    }else{
        return results, nil
    }
}
