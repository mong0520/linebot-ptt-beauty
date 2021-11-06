package controllers

import (
	"errors"
	"log"
	"time"

	"github.com/kkdai/linebot-ptt-beauty/models"
	"github.com/kkdai/linebot-ptt-beauty/utils"
	. "github.com/kkdai/photomgr"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserFavorite struct {
	UserId    string   `json:"user_id" bson:"user_id"`
	Favorites []string `json:"favorites" bson:"favorites"`
}

func GetOne(url string) (result *models.ArticleDocument, err error) {
	ptt := NewPTT()
	post := models.ArticleDocument{}
	post.ArticleID = utils.GetPttIDFromURL(url)
	post.ImageLinks = ptt.GetAllImageAddress(url)
	return &post, nil
}

func Get(page int, perPage int) (results []models.ArticleDocument, err error) {
	var ret []models.ArticleDocument
	ptt := NewPTT()
	count := ptt.ParsePttPageByIndex(page, true)
	for i := 0; i < count && i < perPage; i++ {
		title := ptt.GetPostTitleByIndex(i)
		if utils.CheckTitleWithBeauty(title) {
			post := models.ArticleDocument{}
			url := ptt.GetPostUrlByIndex(i)
			post.ArticleTitle = title
			post.URL = url
			post.ArticleID = utils.GetPttIDFromURL(url)
			post.ImageLinks = ptt.GetAllImageAddress(url)
			like, dis := ptt.GetPostLikeDis(url)
			post.MessageCount.Push = like
			post.MessageCount.Boo = dis
			post.MessageCount.All = like + dis
			ret = append(ret, post)
		}
		// log.Printf("Get article: %s utl= %s obj=%x \n", m.ArticleTitle, m.URL, m)
	}

	return ret, nil
}

func GetRandom(count int, keyword string) (results []models.ArticleDocument, err error) {
	rands := utils.GetRandomIntSet(100, 10)
	ptt := NewPTT()
	pCount := ptt.ParsePttByNumber(101, 0)
	if pCount == 0 {
		return nil, errors.New("NotFound")
	}
	var ret []models.ArticleDocument
	for i := 0; i < count; i++ {
		title := ptt.GetPostTitleByIndex(rands[i])
		if utils.CheckTitleWithBeauty(title) {
			post := models.ArticleDocument{}
			url := ptt.GetPostUrlByIndex(rands[i])
			post.ArticleTitle = title
			post.URL = url
			post.ArticleID = utils.GetPttIDFromURL(url)
			post.ImageLinks = ptt.GetAllImageAddress(url)
			like, dis := ptt.GetPostLikeDis(url)
			post.MessageCount.Push = like
			post.MessageCount.Boo = dis
			post.MessageCount.All = like + dis
			ret = append(ret, post)
			log.Printf("%d th rand =%d title=%s url=%s images(1)=%s \n", i, rands[i], title, url, post.ImageLinks[0])
		}
	}
	return ret, nil
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

func (u *UserFavorite) Add(meta *models.Model) {
	if err := meta.CollectionUserFavorite.Insert(u); err != nil {
		meta.Log.Println(err)
	}
}

func (u *UserFavorite) Get(meta *models.Model) (result *UserFavorite, err error) {
	// meta.Log.Println(u.UserId)
	// query := bson.M{"user_id": u.UserId}
	// if err := meta.CollectionUserFavorite.Find(query).One(&result) ; err != nil{
	//     meta.Log.Println(err)
	//     return nil, err
	// }else{
	return result, nil
	// }
}

func (u *UserFavorite) Update(meta *models.Model) (err error) {
	// meta.Log.Println(u.UserId)
	// query := bson.M{"user_id": u.UserId}
	// if err := meta.CollectionUserFavorite.Update(query, u); err != nil {
	// 	meta.Log.Println(err)
	// 	return err
	// } else {
	return nil
	// }
}
