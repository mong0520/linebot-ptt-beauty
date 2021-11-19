package models

import (
	"encoding/json"
	"log"

	"github.com/go-pg/pg/v10"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	Db  *pg.DB
	Log *log.Logger
}

type MessageCount struct {
	All     int `json:"all" bson:"all"`
	Boo     int `json:"boo" bson:"boo"`
	Count   int `json:"count" bson:"count"`
	Neutral int `json:"neutral" bson:"neutral"`
	Push    int `json:"push" bson:"push"`
}

type ArticleDocument struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ArticleID    string        `json:"article_id" bson:"article_id"`
	ArticleTitle string        `json:"article_title" bson:"article_title"`
	Author       string        `json:"author" bson:"author"`
	Board        string        `json:"board" bson:"board"`
	Content      string        `json:"content" bson:"content"`
	Date         string        `json:"date" bson:"date"`
	IP           string        `json:"ip" bson:"ip"`
	MessageCount MessageCount  `bson:"message_count"`
	Messages     []interface{} `json:"messages" bson:"messages"`
	Timestamp    int           `json:"timestamp" bson:"timestamp"`
	URL          string        `json:"url" bson:"url"`
	ImageLinks   []string      `json:"image_links" bson:"image_links"`
}

func (d *ArticleDocument) GeneralQueryOne(collection *mgo.Collection, query interface{}) (result *ArticleDocument, err error) {
	result = &ArticleDocument{}
	if err := collection.Find(query).One(result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (d *ArticleDocument) GeneralQueryAll(collection *mgo.Collection, query interface{}, sortBy string, count int) (results []ArticleDocument, err error) {
	results = []ArticleDocument{}
	if sortBy == "" {
		if err := collection.Find(query).All(&results); err != nil {
			return nil, err
		} else {
			return results, nil
		}
	} else {
		if err := collection.Find(query).Sort(sortBy).Limit(count).All(&results); err != nil {
			return nil, err
		} else {
			return results, nil
		}
	}

}

func (d *ArticleDocument) ToString() (info string) {
	b, err := json.Marshal(d)
	if err != nil {
		//fmt.Println(err)
		return
	}
	return string(b)
}

// ArticleDocument for sorting.
type AllArticles []ArticleDocument

func (a AllArticles) Len() int           { return len(a) }
func (a AllArticles) Less(i, j int) bool { return a[i].MessageCount.Count > a[j].MessageCount.Count }
func (a AllArticles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
