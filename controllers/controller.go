package controllers

import (
	"errors"
	"fmt"
	"github.com/mong0520/linebot-ptt/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strings"
	"time"
	"sort"
	"github.com/mong0520/linebot-ptt/utils"
)

func GetOne(collection *mgo.Collection, query bson.M) (result *models.ArticleDocument, err error) {
	//query := bson.M{"article_id": "M.1521548086.A.DCA"}
	document := &models.ArticleDocument{}
	result, err = document.GeneralQueryOne(collection, query)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	} else {
		return result, nil
	}
}

func GetAll(collection *mgo.Collection, query bson.M) (results []models.ArticleDocument, err error) {
	document := &models.ArticleDocument{}
	results, err = document.GeneralQueryAll(collection, query, "", -1)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	} else {
		//fmt.Printf("%+v", results)
		return results, nil
	}
}

func GetRandom(collection *mgo.Collection, count int, keyword string) (results []models.ArticleDocument, err error) {
	//document := &models.ArticleDocument{}
	//query := bson.M{"message_count.all": bson.M{"$gt": like}, "ArticleTitle": "/正妹/"}
	//query := bson.M{"ArticleTitle": bson.RegEx{"*", ""}}
	//query := bson.M{"ArticleTitle": bson.RegEx{".+", ""}}
	//query := bson.M{"$sample": bson.M{"size": 10}}
	//log.Println("here")
	query := bson.M{}
	baseline_ts := 1420070400 // 2015年Jan/1/00:00:00 之後
	needRandom := true
	if keyword == "" {
		query = bson.M{"timestamp": bson.M{"$gte": baseline_ts}, "article_title": bson.M{"$regex": bson.RegEx{"^\\[正妹\\].*", ""}}}
	} else {
		query = bson.M{
			"timestamp":     bson.M{"$gte": baseline_ts},
			"article_title": bson.M{"$regex": bson.RegEx{fmt.Sprintf(".*%s.*", strings.ToLower(keyword)), ""}}}
	}

	total, _ := collection.Find(query).Count()
	fmt.Println("total = ", total)
	if total == 0 {
		return nil, errors.New("NotFound")
	}else if total < count{
		count = total
		needRandom = false
	}
	fmt.Println("count = ", count)
	if needRandom{
		randSkip := utils.GetRandomIntSet(total, count)
		//rand.Seed(time.Now().UnixNano())
		for i := 0; i < count; i++ {
			//skip := rand.Intn(total)
			skip := randSkip[i]
			result := &models.ArticleDocument{}
			collection.Find(query).Skip(skip).One(result)
			fmt.Println(skip)
			results = append(results, *result)
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].MessageCount.Push > results[j].MessageCount.Push
		})
	}else{
		document := &models.ArticleDocument{}
		results, err = document.GeneralQueryAll(collection, query, "-message_count.push", count)
	}


	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func GetMostLike(collection *mgo.Collection, count int, timestampOffset int) (results []models.ArticleDocument, err error) {
	document := &models.ArticleDocument{}
	//query := bson.M{"message_count.all": bson.M{"$gt": like}, "ArticleTitle": "/正妹/"}
	//query := bson.M{"ArticleTitle": bson.RegEx{"*", ""}}
	//query := bson.M{"ArticleTitle": bson.RegEx{".+", ""}}
	query := bson.M{}
	if timestampOffset > 0 {
		now := time.Now()
		nowInSec := int(now.Unix())
		start := nowInSec - timestampOffset
		//{"timestamp": {"$gte":  1, "$lt": 9999999999}}
		query = bson.M{"timestamp": bson.M{"$gte": start, "$lt": nowInSec}, "article_title": bson.M{"$regex": bson.RegEx{"^\\[正妹\\].*", ""}}}
	} else {
		query = bson.M{"article_title": bson.M{"$regex": bson.RegEx{"^\\[正妹\\].*", ""}}}
	}
	results, err = document.GeneralQueryAll(collection, query, "-message_count.push", count)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	} else {
		return results, nil
	}
}
