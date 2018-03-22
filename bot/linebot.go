package bot

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/mong0520/linebot-ptt/controllers"
	"github.com/mong0520/linebot-ptt/models"
	"log"
	"mvdan.cc/xurls"
	"net/http"
	"os"
	"strings"
)

var bot *linebot.Client
var meta *models.Model
var maxCountOfCarousel = 10
var defaultImage = "https://s3-ap-northeast-1.amazonaws.com/ottbuilder-neil-test/img/default.png"
var oneDayInSec = 60 * 60 *24
var oneMonthInSec = oneDayInSec * 30
var oneYearInSec = oneMonthInSec * 365

// EventType constants
const (
	ActionDailyHot   string = "本日熱門"
	ActionMonthlyHot string = "近期熱門"
	ActionYearHot    string = "年度熱門"
	ActionRandom     string = "隨機"
	ActionHelp       string = "/show"
)

func InitLineBot(m *models.Model) {
	var err error
	meta = m
	secret := os.Getenv("ChannelSecret")
	token := os.Getenv("ChannelAccessToken")
	bot, err = linebot.New(secret, token)
	if err != nil {
		log.Println(err)
	}
	//log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	//port := os.Getenv("PORT")
	port := "8080"
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			meta.Log.Printf("Receieve Event Type = %s from User [%s] or Room [%s] or Group [%s]\n", event.Type, event.Source.UserID, event.Source.RoomID, event.Source.GroupID)
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				meta.Log.Println("Content = ", message.Text)
				textHander(event, message)
			}
		} else if event.Type == linebot.EventTypePostback {
			meta.Log.Println("got a postback event")
		} else {
			meta.Log.Printf("got a %s event\n", event.Type)
		}
	}
}

func textHander(event *linebot.Event, message *linebot.TextMessage) {
	switch message.Text {
	case ActionDailyHot:
		template := buildCarouseTemplate(ActionDailyHot)
		sendCarouselMessage(event, template)
	case ActionMonthlyHot:
		template := buildCarouseTemplate(ActionMonthlyHot)
		sendCarouselMessage(event, template)
	case ActionYearHot:
		template := buildCarouseTemplate(ActionYearHot)
		sendCarouselMessage(event, template)
	case ActionRandom:
		template := buildCarouseTemplate(ActionRandom)
		sendCarouselMessage(event, template)
	case ActionHelp:
		template := buildButtonTemplate()
		sendButtonMessage(event, template)
	default:
		meta.Log.Println(message.Text)
	}
}

func buildButtonTemplate() (template *linebot.ButtonsTemplate) {
	template = linebot.NewButtonsTemplate("", "表特看看", "你可以試試看...",
		linebot.NewMessageTemplateAction(ActionDailyHot, ActionDailyHot),
		linebot.NewMessageTemplateAction(ActionMonthlyHot, ActionMonthlyHot),
		linebot.NewMessageTemplateAction(ActionYearHot, ActionYearHot),
		linebot.NewMessageTemplateAction(ActionRandom, ActionRandom),
	)
	return template
}

//func buildResponse() (resp string) {
//	results, _ := controllers.GetMostLike(meta.Collection, maxCountOfCarousel)
//	var buffer bytes.Buffer
//	buffer.WriteString("今日熱門表特\n")
//	for _, r := range results {
//		buffer.WriteString(fmt.Sprintf("推文數: {%d}, 標題: {%s}, 網址: {%s}\n", r.MessageCount.All, r.ArticleTitle, r.URL))
//	}
//	resp = buffer.String()
//	log.Println(resp)
//	return resp
//}

//func sendTextMessage(event *linebot.Event, text string) {
//	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
//		log.Println("Send Fail")
//	}
//}

func findImageInContent(content string) (img string) {
	imgs := xurls.Relaxed().FindAllString(content, -1)
	if imgs != nil{
		for _, img := range imgs{
			if strings.HasSuffix(strings.ToLower(img), "jpg"){
				img = strings.Replace(img, "http://", "https://", -1)
				return img
			}
		}
		//meta.Log.Println("try to append jpg in the end")
		img := imgs[0] + ".jpg"
		img = strings.Replace(img, "http://", "https://", -1)
		return img
	}else{
		return defaultImage
	}

}

func buildImgCarouseTemplate(action string) (template *linebot.ImageCarouselTemplate) {
	results := []models.ArticleDocument{}
	switch action {
	case ActionDailyHot:
		results, _ = controllers.GetMostLike(meta.Collection, maxCountOfCarousel, oneDayInSec)
	case ActionRandom:
		results, _ = controllers.GetRandom(meta.Collection, maxCountOfCarousel)
	default:
		results, _ = controllers.GetMostLike(meta.Collection, maxCountOfCarousel, oneMonthInSec)
	}

	columnList := []*linebot.ImageCarouselColumn{}

	for _, result := range results {
		//thumnailUrl := "https://i.imgur.com/mJMuhfP.jpg"
		thumnailUrl := findImageInContent(result.Content)
		//log.Println(thumnailUrl)
		//log.Println(idx, thumnailUrl)
		tmpColumn := linebot.NewImageCarouselColumn(
			thumnailUrl,
			linebot.NewURITemplateAction("點我打開", result.URL),
		)
		columnList = append(columnList, tmpColumn)
	}

	template = linebot.NewImageCarouselTemplate(columnList...)

	return template
}

func buildCarouseTemplate(action string) (template *linebot.CarouselTemplate) {
	results := []models.ArticleDocument{}
	switch action {
	case ActionDailyHot:
		results, _ = controllers.GetMostLike(meta.Collection, maxCountOfCarousel, oneDayInSec)
	case ActionMonthlyHot:
		results, _ = controllers.GetMostLike(meta.Collection, maxCountOfCarousel, oneMonthInSec)
	case ActionYearHot:
		results, _ = controllers.GetMostLike(meta.Collection, maxCountOfCarousel, oneYearInSec)
	case ActionRandom:
		results, _ = controllers.GetRandom(meta.Collection, maxCountOfCarousel)
	default:
		return
	}

	columnList := []*linebot.CarouselColumn{}

	for _, result := range results {
		//meta.Log.Printf("%+v", result)
		//thumnailUrl := "https://c1.sd"
		thumnailUrl := findImageInContent(result.Content)
		//meta.Log.Println(idx, thumnailUrl)
		tmpColumn := linebot.NewCarouselColumn(
			thumnailUrl,
			result.ArticleTitle,
			//fmt.Sprintf("推文數量: %d", result.MessageCount.Push),
			fmt.Sprintf("共有 %d 人推文\n共有 %d 人噓文", result.MessageCount.All, result.MessageCount.Boo),
			linebot.NewURITemplateAction("點我打開", result.URL),
			linebot.NewMessageTemplateAction(ActionDailyHot, ActionDailyHot),
			linebot.NewMessageTemplateAction(ActionRandom, ActionRandom),
		)
		columnList = append(columnList, tmpColumn)
	}

	template = linebot.NewCarouselTemplate(columnList...)

	return template
}

func sendCarouselMessage(event *linebot.Event, template *linebot.CarouselTemplate) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Carousel alt text", template)).Do(); err != nil {
		meta.Log.Println(err)
	}
}

func sendButtonMessage(event *linebot.Event, template *linebot.ButtonsTemplate) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Carousel alt text", template)).Do(); err != nil {
		meta.Log.Println(err)
	}
}

func sendImgCarouseMessage(event *linebot.Event, template *linebot.ImageCarouselTemplate) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Carousel alt text", template)).Do(); err != nil {
		meta.Log.Println(err)
	}
}
