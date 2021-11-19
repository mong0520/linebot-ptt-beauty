package controllers

import (
	"errors"
	"log"
	"sort"

	"github.com/kkdai/linebot-ptt-beauty/models"
	"github.com/kkdai/linebot-ptt-beauty/utils"
	. "github.com/kkdai/photomgr"
)

type UserFavorite struct {
	Id        int64    `bson:"_id"`
	UserId    string   `json:"user_id" bson:"user_id"`
	Favorites []string `json:"favorites" bson:"favorites"`
}

func GetOne(url string) (result *models.ArticleDocument, err error) {
	ptt := NewPTT()
	post := models.ArticleDocument{}
	post.URL = url
	post.ArticleID = utils.GetPttIDFromURL(url)
	post.ArticleTitle = ptt.GetUrlTitle(url)
	post.ImageLinks = ptt.GetAllImageAddress(url)
	like, dis := ptt.GetPostLikeDis(url)
	post.MessageCount.Push = like
	post.MessageCount.Boo = dis
	post.MessageCount.All = like + dis
	post.MessageCount.Count = like - dis
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
			post.MessageCount.Count = ptt.GetPostStarByIndex(i)
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
			post.MessageCount.Count = ptt.GetPostStarByIndex(i)
			ret = append(ret, post)
			// log.Printf("%d th rand =%d title=%s url=%s images(1)=%s \n", i, rands[i], title, url, post.ImageLinks[0])
		}
	}
	return ret, nil
}

func GetMostLike(count int, timestampOffset int) (results []models.ArticleDocument, err error) {
	// if timestampOffset > 0 {
	// 	now := time.Now()
	// 	nowInSec := int(now.Unix())
	// 	start := nowInSec - timestampOffset
	// 	//{"timestamp": {"$gte":  1, "$lt": 9999999999}}
	// 	query = bson.M{"timestamp": bson.M{"$gte": start, "$lt": nowInSec}, "article_title": bson.M{"$regex": bson.RegEx{"^\\[正妹\\].*", ""}}}
	// } else {
	// 	query = bson.M{"article_title": bson.M{"$regex": bson.RegEx{"^\\[正妹\\].*", ""}}}
	// }

	ptt := NewPTT()
	pCount := ptt.ParsePttByNumber(20, 0)
	if pCount == 0 {
		return nil, errors.New("NotFound")
	}

	var ret []models.ArticleDocument
	for i := 0; i < count && i < pCount; i++ {
		title := ptt.GetPostTitleByIndex(i)
		post := models.ArticleDocument{}
		url := ptt.GetPostUrlByIndex(i)
		like, dis := ptt.GetPostLikeDis(url)
		post.MessageCount.Push = like
		post.MessageCount.Boo = dis
		post.MessageCount.All = like + dis
		post.MessageCount.Count = ptt.GetPostStarByIndex(i)
		post.ArticleTitle = title
		post.URL = url
		post.ArticleID = utils.GetPttIDFromURL(url)
		post.ImageLinks = ptt.GetAllImageAddress(url)
		ret = append(ret, post)
		// log.Printf("%d  stars=%d  title=%s\n", i, post.MessageCount.Count, title)
	}
	sort.Sort(models.AllArticles(ret))
	// for i := 0; i < count; i++ {
	// 	log.Printf("%d  stars=%d  title=%s\n", i, ret[i].MessageCount.Count, ret[i].ArticleTitle)
	// }
	return ret[0:count], nil
}

func (u *UserFavorite) Add(meta *models.Model) {
	_, err := meta.Db.Model(u).Insert()
	if err != nil {
		meta.Log.Println(err)
	}
}

func (u *UserFavorite) Get(meta *models.Model) (result *UserFavorite, err error) {
	log.Println("***Get Fav uUID=", u.UserId)
	userFav := UserFavorite{}
	err = meta.Db.Model(&userFav).
		Where("user_id = ?", u.UserId).
		Select()
	if err != nil {
		meta.Log.Println(err)
		return nil, err
	}
	meta.Log.Println("UserFavorite DB result= ", userFav)
	return &userFav, nil
}

func (u *UserFavorite) Update(meta *models.Model) (err error) {
	log.Println("***Update Fav User=", u)

	_, err = meta.Db.Model(u).
		Set("favorites = ?", u.Favorites).
		Where("user_id = ?", u.UserId).
		Update()
	if err != nil {
		meta.Log.Println(err)
	}
	return nil
}

func (u *UserFavorite) ShowAll(meta *models.Model) (err error) {
	log.Println("***ShowAll  User -->")
	users := []UserFavorite{}
	err = meta.Db.Model(&users).Select()
	if err != nil {
		log.Println(err)
	}
	log.Println("***Show all users =", users)
	return nil
}
