package bots

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/mong0520/linebot-ptt-beauty/controllers"
	"github.com/mong0520/linebot-ptt-beauty/models"
	"github.com/mong0520/linebot-ptt-beauty/utils"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var bot *linebot.Client
var meta *models.Model
var maxCountOfCarousel = 10
var defaultImage = "https://i.imgur.com/WAnWk7K.png"
var defaultThumbnail = "https://i.imgur.com/StcRAPB.png"
var oneDayInSec = 60 * 60 * 24
var oneMonthInSec = oneDayInSec * 30
var oneYearInSec = oneMonthInSec * 365
var SSLCertPath = "/etc/dehydrated/certs/nt1.me/fullchain.pem"
var SSLPrivateKeyPath = "/etc/dehydrated/certs/nt1.me/privkey.pem"

// EventType constants
const (
	DefaultTitle string = "💋表特看看"

	// 應該把 action 和 lable 分開
	ActionQuery       string = "一般查詢"
	ActionNewest      string = "🎊 最新表特"
	ActionDailyHot    string = "📈 本日熱門"
	ActionMonthlyHot  string = "🔥 近期熱門" //改成近期隨機, 先選出100個，然後隨機吐10筆
	ActionYearHot     string = "🏆 年度熱門"
	ActionRandom      string = "👩 隨機十連抽"
	ActionAddFavorite string = "加入最愛"
	ActionClick       string = "👉 點我打開"
	ActionHelp        string = "表特選單"
	ActionAllImage    string = "👁️ 預覽圖片"
	ActonShowFav      string = "❤️ 我的最愛"

	ModeHttp  string = "http"
	ModeHttps string = "https"
	AltText   string = "正妹只在手機上"
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
	port := os.Getenv("PORT")
	//port := "8080"
	addr := fmt.Sprintf(":%s", port)
	runMode := os.Getenv("RUNMODE")
	m.Log.Printf("Run Mode = %s\n", runMode)
	if strings.ToLower(runMode) == ModeHttps {
		m.Log.Printf("Secure listen on %s with \n", addr)
		http.ListenAndServeTLS(addr, SSLCertPath, SSLPrivateKeyPath, nil)
	} else {
		m.Log.Printf("Listen on %s\n", addr)
		http.ListenAndServe(addr, nil)
	}
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
			userDisplayName := getUserNameById(event.Source.UserID)
			meta.Log.Printf("Receieve Event Type = %s from User [%s](%s), or Room [%s] or Group [%s]\n",
				event.Type, userDisplayName, event.Source.UserID, event.Source.RoomID, event.Source.GroupID)

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				meta.Log.Println("Text = ", message.Text)
				textHander(event, message.Text)
			default:
				meta.Log.Println("Unimplemented handler for event type ", event.Type)
			}
		} else if event.Type == linebot.EventTypePostback {
			meta.Log.Println("got a postback event")
			meta.Log.Println(event.Postback.Data)
			postbackHandler(event)

		} else {
			meta.Log.Printf("got a %s event\n", event.Type)
		}
	}
}

func actionHandler(event *linebot.Event, action string, values url.Values) {
	switch action {
	case ActionNewest:
		actionNewest(event, values)
	case ActionAllImage:
		actionAllImage(event, values)
	case ActionQuery, ActionRandom:
		actionGeneral(event, action, values)
	case ActionAddFavorite:
		actinoAddFavorite(event, action, values)
	case ActonShowFav:
		actionShowFavorite(event, action, values)
	default:
		meta.Log.Println("Unimplement action handler", action)
	}
}

func actinoAddFavorite(event *linebot.Event, action string, values url.Values) {
	toggleMessage := ""
	userId := values.Get("user_id")
	newFavoriteArticle := values.Get("article_id")
	userFavorite := &controllers.UserFavorite{
		UserId:    userId,
		Favorites: []string{newFavoriteArticle},
	}
	latestFavArticles := []string{}
	if record, err := userFavorite.Get(meta); err != nil {
		meta.Log.Println("User data is not created, create a new one")
		userFavorite.Add(meta)
		latestFavArticles = append(latestFavArticles, newFavoriteArticle)
	} else {
		meta.Log.Println("Record found, update it", record)
		oldRecords := record.Favorites
		if exist, idx := utils.InArray(newFavoriteArticle, oldRecords); exist == true {
			meta.Log.Println(newFavoriteArticle, "已存在，移除")
			oldRecords = utils.RemoveStringItem(oldRecords, idx)
			toggleMessage = "已從最愛中移除"
		} else {
			oldRecords = append(oldRecords, newFavoriteArticle)
			toggleMessage = "已新增至最愛"
		}
		latestFavArticles = oldRecords
		userFavorite.Favorites = oldRecords
		userFavorite.Update(meta)
	}
	sendTextMessage(event, toggleMessage)
}

func actionShowFavorite(event *linebot.Event, action string, values url.Values) {
	userFavorite := &controllers.UserFavorite{
		UserId:    values.Get("user_id"),
		Favorites: []string{},
	}
	userData, _ := userFavorite.Get(meta)
	favDocuments := []models.ArticleDocument{}
	for idx := len(userData.Favorites)-1; idx >= 0; idx-- {
		favArticleId := userData.Favorites[idx]
		query := bson.M{"article_id": favArticleId}
		tmpRecord, _ := controllers.GetOne(meta.Collection, query)
		favDocuments = append(favDocuments, *tmpRecord)
	}
	if len(favDocuments) == 0 {
		sendTextMessage(event, "尚無最愛")
	} else {
		template := getCarouseTemplate(event.Source.UserID, favDocuments)
		sendCarouselMessage(event, template, "最愛照片已送達")
	}
}

func actionGeneral(event *linebot.Event, action string, values url.Values) {
	meta.Log.Println("Enter actionGeneral, action = ", action)
	meta.Log.Println("Enter actionGeneral, values = ", values)
	records := []models.ArticleDocument{}
	label := ""
	switch action {
	case ActionQuery:
		//meta.Log.Println(values.Get("period"))
		tsOffset, _ := strconv.Atoi(values.Get("period"))
		meta.Log.Println("timestampe off set = ", tsOffset)
		records, _ = controllers.GetMostLike(meta.Collection, maxCountOfCarousel, tsOffset)
		label = "已幫您查詢到一些照片~"
	case ActionRandom:
		records, _ = controllers.GetRandom(meta.Collection, maxCountOfCarousel, "")
		label = "隨機表特已送到囉"
	default:
		return
	}
	template := getCarouseTemplate(event.Source.UserID, records)
	if template != nil {
		sendCarouselMessage(event, template, label)
	}

}

func actionAllImage(event *linebot.Event, values url.Values) {
	if articleId := values.Get("article_id"); articleId != "" {
		query := bson.M{"article_id": articleId}
		result, _ := controllers.GetOne(meta.Collection, query)
		template := getImgCarousTemplate(result)
		sendImgCarouseMessage(event, template)
	} else {
		meta.Log.Println("Unable to get article id", values)
	}
}

func actionNewest(event *linebot.Event, values url.Values) {
	columnCount := 9
	if currentPage, err := strconv.Atoi(values.Get("page")); err != nil {
		meta.Log.Println("Unable to parse parameters", values)
	} else {
		records, _ := controllers.Get(meta.Collection, currentPage, columnCount)
		template := getCarouseTemplate(event.Source.UserID, records)

		if template == nil {
			meta.Log.Println("Unable to get template", values)
			return
		}

		// append next page column
		previousPage := currentPage - 1
		if previousPage < 0 {
			previousPage = 0
		}
		nextPage := currentPage + 1
		previousData := fmt.Sprintf("action=%s&page=%d", ActionNewest, previousPage)
		nextData := fmt.Sprintf("action=%s&page=%d", ActionNewest, nextPage)
		previousText := fmt.Sprintf("上一頁 %d", previousPage)
		nextText := fmt.Sprintf("下一頁 %d", nextPage)
		tmpColumn := linebot.NewCarouselColumn(
			defaultThumbnail,
			DefaultTitle,
			"繼續看？",
			linebot.NewMessageTemplateAction(ActionHelp, ActionHelp),
			linebot.NewPostbackTemplateAction(previousText, previousData, "", ""),
			linebot.NewPostbackTemplateAction(nextText, nextData, "", ""),
		)
		template.Columns = append(template.Columns, tmpColumn)

		sendCarouselMessage(event, template, "熱騰騰的最新照片送到了!")
	}
}

func getCarouseTemplate(userId string, records []models.ArticleDocument) (template *linebot.CarouselTemplate) {
	if len(records) == 0 {
		return nil
	}

	columnList := []*linebot.CarouselColumn{}
	userFavorite := &controllers.UserFavorite{
		UserId:    userId,
		Favorites: []string{},
	}
	userData, _ := userFavorite.Get(meta)
	favLabel := ""

	for _, result := range records {
		if exist, _ := utils.InArray(result.ArticleID, userData.Favorites); exist == true {
			favLabel = "❤️ 移除最愛"
		} else {
			favLabel = "💛 加入最愛"
		}
		thumnailUrl := defaultImage
		imgUrlCounts := len(result.ImageLinks)
		lable := fmt.Sprintf("%s (%d)", ActionAllImage, imgUrlCounts)
		title := result.ArticleTitle
		postBackData := fmt.Sprintf("action=%s&article_id=%s&page=0", ActionAllImage, result.ArticleID)
		text := fmt.Sprintf("%d 😍\t%d 😡", result.MessageCount.Push, result.MessageCount.Boo)

		if imgUrlCounts > 0 {
			thumnailUrl = result.ImageLinks[0]
		}

		// Title's hard limit by Line
		if len(title) >= 40 {
			title = title[0:38]
		}
		//meta.Log.Println("===============", idx)
		//meta.Log.Println("Thumbnail Url = ", thumnailUrl)
		//meta.Log.Println("Title = ", title)
		//meta.Log.Println("Text = ", text)
		//meta.Log.Println("URL = ", result.URL)
		//meta.Log.Println("===============", idx)
		//dataRandom := fmt.Sprintf("action=%s", ActionRandom)
		dataAddFavorite := fmt.Sprintf("action=%s&user_id=%s&article_id=%s",
			ActionAddFavorite, userId, result.ArticleID)
		tmpColumn := linebot.NewCarouselColumn(
			thumnailUrl,
			title,
			text,
			linebot.NewURITemplateAction(ActionClick, result.URL),
			linebot.NewPostbackTemplateAction(lable, postBackData, "", ""),
			//linebot.NewPostbackTemplateAction(ActionRandom, dataRandom, "", ""),
			linebot.NewPostbackTemplateAction(favLabel, dataAddFavorite, "", ""),
		)
		columnList = append(columnList, tmpColumn)
	}
	template = linebot.NewCarouselTemplate(columnList...)
	return template
}

func postbackHandler(event *linebot.Event) {
	m, _ := url.ParseQuery(event.Postback.Data)
	action := m.Get("action")
	meta.Log.Println("Action = ", action)
	actionHandler(event, action, m)
}

func getUserNameById(userId string) (userDisplayName string) {
	res, err := bot.GetProfile(userId).Do()
	if err != nil {
		userDisplayName = "Unknown"
	} else {
		userDisplayName = res.DisplayName
	}
	return userDisplayName
}

func textHander(event *linebot.Event, message string) {
	userFavorite := &controllers.UserFavorite{
		UserId:    event.Source.UserID,
		Favorites: []string{},
	}
	if _, err := userFavorite.Get(meta); err != nil {
		meta.Log.Println("User data is not created, create a new one")
		userFavorite.Add(meta)
	}
	switch message {
	case ActionHelp:
		template := getMenuButtonTemplateV2(event, DefaultTitle)
		sendCarouselMessage(event, template, "我能為您做什麼？")
	case ActionRandom:
		records, _ := controllers.GetRandom(meta.Collection, maxCountOfCarousel, "")
		template := getCarouseTemplate(event.Source.UserID, records)
		sendCarouselMessage(event, template, "隨機表特已送到囉")
	case ActionNewest:
		values := url.Values{}
		values.Set("period", fmt.Sprintf("%d", oneDayInSec))
		values.Set("page", "0")
		actionNewest(event, values)
    case ActonShowFav:
        values := url.Values{}
        values.Set("user_id", event.Source.UserID)
        actionShowFavorite(event, "", values)
	default:
		if event.Source.UserID != "" && event.Source.GroupID == "" && event.Source.RoomID == "" {
			records, _ := controllers.GetRandom(meta.Collection, maxCountOfCarousel, message)
			if records != nil && len(records) > 0 {
				template := getCarouseTemplate(event.Source.UserID, records)
				sendCarouselMessage(event, template, "隨機表特已送到囉")
			} else {
				template := getMenuButtonTemplateV2(event, DefaultTitle)
				sendCarouselMessage(event, template, "我能為您做什麼？")
			}
		}
	}
}

func getMenuButtonTemplateV2(event *linebot.Event, title string) (template *linebot.CarouselTemplate) {
	columnList := []*linebot.CarouselColumn{}
	dataNewlest := fmt.Sprintf("action=%s&page=0", ActionNewest)
	dataRandom := fmt.Sprintf("action=%s", ActionRandom)
	dataQuery := fmt.Sprintf("action=%s", ActionQuery)
	dataShowFav := fmt.Sprintf("action=%s&user_id=%s", ActonShowFav, event.Source.UserID)

	menu1 := linebot.NewCarouselColumn(
		defaultThumbnail,
		title,
		"你可以試試看以下選項，或直接輸入關鍵字查詢",
		linebot.NewPostbackTemplateAction(ActionNewest, dataNewlest, "", ""),
		linebot.NewPostbackTemplateAction(ActionRandom, dataRandom, "", ""),
		linebot.NewPostbackTemplateAction(ActonShowFav, dataShowFav, "", ""),
	)
	menu2 := linebot.NewCarouselColumn(
		defaultThumbnail,
		title,
		"你可以試試看以下選項，或直接輸入關鍵字查詢",
		linebot.NewPostbackTemplateAction(ActionDailyHot, dataQuery+"&period="+fmt.Sprintf("%d", oneDayInSec), "", ""),
		linebot.NewPostbackTemplateAction(ActionMonthlyHot, dataQuery+"&period="+fmt.Sprintf("%d", oneMonthInSec), "", ""),
		linebot.NewPostbackTemplateAction(ActionYearHot, dataQuery+"&period="+fmt.Sprintf("%d", oneYearInSec), "", ""),
	)
	columnList = append(columnList, menu1, menu2)
	template = linebot.NewCarouselTemplate(columnList...)
	return template
}

//func getMenuButtonTemplate(event *linebot.Event, title string) (template *linebot.ButtonsTemplate) {
//	dataNewlest := fmt.Sprintf("action=%s&page=0", ActionNewest)
//	dataRandom := fmt.Sprintf("action=%s", ActionRandom)
//	dataQuery := fmt.Sprintf("action=%s", ActionQuery)
//	dataShowFav := fmt.Sprintf("action=%s&user_id=%s", ActonShowFav, event.Source.UserID)
//	template = linebot.NewButtonsTemplate(defaultThumbnail, title, "你可以試試看以下選項，或直接輸入關鍵字查詢",
//		linebot.NewPostbackTemplateAction(ActionNewest, dataNewlest, "", ""),
//		linebot.NewPostbackTemplateAction(ActionDailyHot, dataQuery+"&period="+fmt.Sprintf("%d", oneDayInSec), "", ""),
//		linebot.NewPostbackTemplateAction(ActonShowFav, dataShowFav, "", ""),
//		//linebot.NewPostbackTemplateAction(ActionMonthlyHot, dataQuery+"&period="+fmt.Sprintf("%d", oneMonthInSec), "", ""),
//		//linebot.NewPostbackTemplateAction(ActionYearHot, dataQuery + "&period="+fmt.Sprintf("%d", oneYearInSec), "", ""),
//		linebot.NewPostbackTemplateAction(ActionRandom, dataRandom, "", ""),
//	)
//	return template
//}

func sendTextMessage(event *linebot.Event, text string) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
		log.Println("Send Fail")
	}
}

func getImgCarousTemplate(record *models.ArticleDocument) (template *linebot.ImageCarouselTemplate) {
	urls := record.ImageLinks
	columnList := []*linebot.ImageCarouselColumn{}
	if len(urls) > 10 {
		urls = urls[0:10]
	}
	for _, url := range urls {
		tmpColumn := linebot.NewImageCarouselColumn(
			url,
			linebot.NewURITemplateAction(ActionClick, record.URL),
		)
		columnList = append(columnList, tmpColumn)
	}
	template = linebot.NewImageCarouselTemplate(columnList...)
	return template
}

func sendCarouselMessage(event *linebot.Event, template *linebot.CarouselTemplate, altText string) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage(altText, template)).Do(); err != nil {
		meta.Log.Println(err)
	}
}

//func sendButtonMessage(event *linebot.Event, template *linebot.ButtonsTemplate) {
//	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage(AltText, template)).Do(); err != nil {
//		meta.Log.Println(err)
//	}
//}

func sendImgCarouseMessage(event *linebot.Event, template *linebot.ImageCarouselTemplate) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("預覽圖片已送達", template)).Do(); err != nil {
		meta.Log.Println(err)
	}
}
