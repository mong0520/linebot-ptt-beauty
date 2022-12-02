package bots

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"mvdan.cc/xurls/v2"

	"github.com/kkdai/linebot-ptt-beauty/controllers"
	"github.com/kkdai/linebot-ptt-beauty/models"
	"github.com/kkdai/linebot-ptt-beauty/utils"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client
var meta *models.Model
var maxCountOfCarousel = 10
var defaultImage = "https://i.imgur.com/WAnWk7K.png"
var defaultThumbnail = "https://i.imgur.com/StcRAPB.png"

const (
	//DefaultTitle : for caresoul title.
	DefaultTitle string = "💋表特看看"

	ActionQuery       string = "一般查詢"
	ActionNewest      string = "🎊 最新表特"
	ActionDailyHot    string = "📈 20篇內熱門"
	ActionMonthlyHot  string = "🔥 60篇內熱門"
	ActionYearHot     string = "🏆 100篇內熱門"
	ActionRandom      string = "👩 隨機十連抽"
	ActionAddFavorite string = "加入最愛"
	ActionClick       string = "👉 點我打開"
	ActionHelp        string = "表特選單"
	ActionAllImage    string = "👁️ 預覽圖片"
	ActonShowFav      string = "❤️ 我的最愛"
	ModeHTTP          string = "http"
	ModeHTTPS         string = "https"
	AltText           string = "正妹只在手機上"
)

// InitLineBot: init LINE bot
func InitLineBot(m *models.Model, runMode string, sslCertPath string, sslPKeyPath string) {

	var err error
	meta = m
	secret := os.Getenv("ChannelSecret")
	token := os.Getenv("ChannelAccessToken")
	bot, err = linebot.New(secret, token)
	if err != nil {
		log.Println(err)
	}
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/health", healthHandler)
	port := os.Getenv("PORT")

	addr := fmt.Sprintf(":%s", port)
	m.Log.Printf("Run Mode = %s\n", runMode)
	if strings.ToLower(runMode) == ModeHTTPS {
		m.Log.Printf("Secure listen on %s with \n", addr)
		err := http.ListenAndServeTLS(addr, sslCertPath, sslPKeyPath, nil)
		if err != nil {
			m.Log.Panic(err)
		}
	} else {
		m.Log.Printf("Listen on %s\n", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			m.Log.Panic(err)
		}
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	meta.Log.Println("enter callback hander")
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
				if message.Text == "showall" {
					log.Println("get show all user OP--->")
					userFavorite := &controllers.UserFavorite{
						UserId:    event.Source.UserID,
						Favorites: []string{},
					}
					userFavorite.ShowAll(meta)
					sendTextMessage(event, "Already show all user DB OP.")
					return
				}
				if strings.Contains(message.Text, "www.ptt.cc/bbs/Beauty") {
					values := url.Values{}
					values.Set("user_id", event.Source.UserID)
					rxRelaxed := xurls.Relaxed()
					values.Set("url", rxRelaxed.FindString(message.Text))
					actinoAddFavorite(event, "", values)
					return
				}

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
	case ActionQuery, ActionDailyHot, ActionMonthlyHot, ActionYearHot:
		actionMostLike(event, action, values)
	case ActionRandom:
		actionRandom(event, values)
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
	newFavoriteArticle := values.Get("url")
	userFavorite := &controllers.UserFavorite{
		UserId:    userId,
		Favorites: []string{newFavoriteArticle},
	}
	log.Println("Add Fav UID", userFavorite.UserId, " Fav[]=", userFavorite.Favorites)
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
	columnCount := 9
	userId := values.Get("user_id")
	userFavorite := &controllers.UserFavorite{
		UserId:    userId,
		Favorites: []string{},
	}

	if currentPage, err := strconv.Atoi(values.Get("page")); err != nil {
		meta.Log.Println("Unable to parse parameters", values)
	} else {
		userData, _ := userFavorite.Get(meta)

		// 	// reverse slice
		for i := len(userData.Favorites)/2 - 1; i >= 0; i-- {
			opp := len(userData.Favorites) - 1 - i
			userData.Favorites[i], userData.Favorites[opp] = userData.Favorites[opp], userData.Favorites[i]
		}

		startIdx := currentPage * columnCount
		endIdx := startIdx + columnCount
		lastPage := false
		if endIdx > len(userData.Favorites)-1 || startIdx > endIdx {
			endIdx = len(userData.Favorites)
			lastPage = true
		}

		// log.Println("Start Index", startIdx)
		// log.Println("End Index", endIdx)
		// log.Println("Total Length", len(userData.Favorites))

		favDocuments := []models.ArticleDocument{}
		favs := userData.Favorites[startIdx:endIdx]
		log.Println(favs)

		for i := startIdx; i < endIdx; i++ {
			url := userData.Favorites[i]
			tmpRecord, _ := controllers.GetOne(url)
			// log.Printf("Favorites[%d] url=%s title=%s \n", i, url, tmpRecord.ArticleTitle)
			favDocuments = append(favDocuments, *tmpRecord)
		}

		// append next page column
		previousPage := currentPage - 1
		if previousPage < 0 {
			previousPage = 0
		}
		nextPage := currentPage + 1
		previousData := fmt.Sprintf("action=%s&page=%d&user_id=%s", ActonShowFav, previousPage, userId)
		nextData := fmt.Sprintf("action=%s&page=%d&user_id=%s", ActonShowFav, nextPage, userId)
		previousText := fmt.Sprintf("上一頁 %d", previousPage)
		nextText := fmt.Sprintf("下一頁 %d", nextPage)
		if lastPage == true {
			nextData = "--"
			nextText = "--"
		}

		tmpColumn := linebot.NewCarouselColumn(
			defaultThumbnail,
			DefaultTitle,
			"繼續看？",
			linebot.NewMessageAction(ActionHelp, ActionHelp),
			linebot.NewPostbackAction(previousText, previousData, "", "", "", ""),
			linebot.NewPostbackAction(nextText, nextData, "", "", "", ""),
		)

		template := getCarouseTemplate(event.Source.UserID, favDocuments)
		template.Columns = append(template.Columns, tmpColumn)
		sendCarouselMessage(event, template, "最愛照片已送達")
	}
}

func actionRandom(event *linebot.Event, values url.Values) {
	var label string
	records, _ := controllers.GetRandom(maxCountOfCarousel, "")
	label = "隨機表特已送到囉"
	template := getCarouseTemplate(event.Source.UserID, records)
	if template != nil {
		sendCarouselMessage(event, template, label)
	}
}

func actionMostLike(event *linebot.Event, action string, values url.Values) {
	period, _ := strconv.Atoi(values.Get("period"))
	records, _ := controllers.GetMostLike(period, maxCountOfCarousel)
	label := "已幫您查詢到一些照片~"

	template := getCarouseTemplate(event.Source.UserID, records)
	if template != nil {
		sendCarouselMessage(event, template, label)
	}
}

func actionAllImage(event *linebot.Event, values url.Values) {
	if url := values.Get("url"); url != "" {
		result, _ := controllers.GetOne(url)
		template := getImgCarousTemplate(result, values)
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
		records, _ := controllers.Get(currentPage, columnCount)
		// in case page 0 is no girls.
		if len(records) == 0 {
			currentPage++
			records, _ = controllers.Get(currentPage, columnCount)
		}

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
			linebot.NewMessageAction(ActionHelp, ActionHelp),
			linebot.NewPostbackAction(previousText, previousData, "", "", "", ""),
			linebot.NewPostbackAction(nextText, nextData, "", "", "", ""),
		)
		template.Columns = append(template.Columns, tmpColumn)

		sendCarouselMessage(event, template, "熱騰騰的最新照片送到了!")
	}
}

func getCarouseTemplate(userId string, records []models.ArticleDocument) (template *linebot.CarouselTemplate) {
	if len(records) == 0 {
		log.Println("err1")
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
		if exist, _ := utils.InArray(result.URL, userData.Favorites); exist == true {
			favLabel = "❤️ 移除最愛"
		} else {
			favLabel = "💛 加入最愛"
		}
		thumnailUrl := defaultImage
		imgUrlCounts := len(result.ImageLinks)
		lable := fmt.Sprintf("%s (%d)", ActionAllImage, imgUrlCounts)
		title := result.ArticleTitle
		postBackData := fmt.Sprintf("action=%s&page=0&url=%s", ActionAllImage, result.URL)
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
		dataAddFavorite := fmt.Sprintf("action=%s&user_id=%s&url=%s",
			ActionAddFavorite, userId, result.URL)
		tmpColumn := linebot.NewCarouselColumn(
			thumnailUrl,
			title,
			text,
			linebot.NewURIAction(ActionClick, result.URL),
			linebot.NewPostbackAction(lable, postBackData, "", "", "", ""),
			//linebot.NewPostbackAction(ActionRandom, dataRandom, "", ""),
			linebot.NewPostbackAction(favLabel, dataAddFavorite, "", "", "", ""),
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
	log.Println("txMSG=", message)
	switch message {
	case "Menu", "menu", "Help", "help", ActionHelp:
		template := getMenuButtonTemplateV2(event, DefaultTitle)
		sendCarouselMessage(event, template, "我能為您做什麼？")
	case ActionRandom:
		records, _ := controllers.GetRandom(maxCountOfCarousel, "")
		template := getCarouseTemplate(event.Source.UserID, records)
		sendCarouselMessage(event, template, "隨機表特已送到囉")
	case ActionNewest:
		values := url.Values{}
		values.Set("page", "0")
		actionNewest(event, values)
	case ActonShowFav:
		values := url.Values{}
		values.Set("user_id", event.Source.UserID)
		values.Set("page", "0")
		actionShowFavorite(event, "", values)
	default:
		template := getMenuButtonTemplateV2(event, DefaultTitle)
		sendCarouselMessage(event, template, "我能為您做什麼？")
		return
	}

	if event.Source.UserID != "" && event.Source.GroupID == "" && event.Source.RoomID == "" {
		records, _ := controllers.GetRandom(maxCountOfCarousel, message)
		if len(records) > 0 {
			template := getCarouseTemplate(event.Source.UserID, records)
			sendCarouselMessage(event, template, "隨機表特已送到囉")
		} else {
			template := getMenuButtonTemplateV2(event, DefaultTitle)
			sendCarouselMessage(event, template, "我能為您做什麼？")
		}
	}
}

func getMenuButtonTemplateV2(event *linebot.Event, title string) (template *linebot.CarouselTemplate) {
	columnList := []*linebot.CarouselColumn{}
	dataNewlest := fmt.Sprintf("action=%s&page=0", ActionNewest)
	dataRandom := fmt.Sprintf("action=%s", ActionRandom)
	dataQuery := fmt.Sprintf("action=%s", ActionQuery)
	dataShowFav := fmt.Sprintf("action=%s&user_id=%s&page=0", ActonShowFav, event.Source.UserID)

	menu1 := linebot.NewCarouselColumn(
		defaultThumbnail,
		title,
		"你可以試試看以下選項，或直接輸入關鍵字查詢",
		linebot.NewPostbackAction(ActionNewest, dataNewlest, "", "", "", ""),
		linebot.NewPostbackAction(ActionRandom, dataRandom, "", "", "", ""),
		linebot.NewPostbackAction(ActonShowFav, dataShowFav, "", "", "", ""),
	)
	menu2 := linebot.NewCarouselColumn(
		defaultThumbnail,
		title,
		"你可以試試看以下選項，或直接輸入關鍵字查詢",
		linebot.NewPostbackAction(ActionDailyHot, dataQuery+"&period=20", "", "", "", ""),
		linebot.NewPostbackAction(ActionMonthlyHot, dataQuery+"&period=60", "", "", "", ""),
		linebot.NewPostbackAction(ActionYearHot, dataQuery+"&period=100", "", "", "", ""),
	)
	columnList = append(columnList, menu1, menu2)
	template = linebot.NewCarouselTemplate(columnList...)
	return template
}

func sendTextMessage(event *linebot.Event, text string) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
		log.Println("Send Fail")
	}
}

func getImgCarousTemplate(record *models.ArticleDocument, values url.Values) (template *linebot.ImageCarouselTemplate) {
	urls := record.ImageLinks
	columnList := []*linebot.ImageCarouselColumn{}
	targetUrl := values.Get("url")
	log.Println("fix img url=", targetUrl)
	page, _ := strconv.Atoi(values.Get("page"))
	startIdx := page * 9
	endIdx := startIdx + 9
	lastPage := false
	if endIdx >= len(urls)-1 {
		endIdx = len(urls)
		lastPage = true
	}
	urls = urls[startIdx:endIdx]

	for _, url := range urls {
		tmpColumn := linebot.NewImageCarouselColumn(
			url,
			linebot.NewURIAction(ActionClick, url),
		)
		columnList = append(columnList, tmpColumn)
	}
	if !lastPage {
		postBackData := fmt.Sprintf("action=%s&page=%d&url=%s", ActionAllImage, page+1, targetUrl)
		tmpColumn := linebot.NewImageCarouselColumn(
			defaultImage,
			linebot.NewPostbackAction("下一頁", postBackData, "", "", "", ""),
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

func sendImgCarouseMessage(event *linebot.Event, template *linebot.ImageCarouselTemplate) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("預覽圖片已送達", template)).Do(); err != nil {
		meta.Log.Println(err)
	}
}
