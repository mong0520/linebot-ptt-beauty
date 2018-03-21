package bot

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/mong0520/linebot-ptt/controllers"
	"github.com/mong0520/linebot-ptt/models"
	"log"
	"net/http"
    "bytes"
    "os"
    "mvdan.cc/xurls"
)

var bot *linebot.Client
var meta *models.Model

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
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Println("Receieve message " + message.Text)
				textHander(event, message)
			}
		}
	}
}

func textHander(event *linebot.Event, message *linebot.TextMessage) {
	switch message.Text {
	case "show me":
        //response := buildResponse()
        //sendTextMessage(event, response)

        template := buildCarouseTemplate()
        sendCarouselMessage(event, template)
	}
	//if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
	//
	//}
}

func buildResponse() (resp string){
	results, _ := controllers.GetMostLike(meta.Collection, 5)
    var buffer bytes.Buffer
    buffer.WriteString("今日熱門表特\n")
	for _, r := range results {
        buffer.WriteString(fmt.Sprintf("推文數: {%d}, 標題: {%s}, 網址: {%s}\n", r.MessageCount.All, r.ArticleTitle, r.URL))
	}
    resp = buffer.String()
    log.Println(resp)
	return resp
}


func sendTextMessage(event *linebot.Event, text string){
    if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
        log.Println("Send Fail")
    }
}

func findImageInContent(content string)(img string){
	img = xurls.Relaxed().FindString(content)
	return img
}

func buildCarouseTemplate()(template *linebot.CarouselTemplate){
	results, _ := controllers.GetMostLike(meta.Collection, 5)

	columnList := []linebot.CarouselColumn{}

	for _, result := range results{
		thumnailUrl := findImageInContent(result.Content)
		fmt.Println(thumnailUrl)
		tmpColumn := linebot.NewCarouselColumn(
			thumnailUrl,
			result.ArticleTitle,
			fmt.Sprintf("推文數量: %d", result.MessageCount.Push),
			linebot.NewURITemplateAction("點我打開", result.URL),
			//linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", "", ""),
		)
		columnList = append(columnList, *tmpColumn)
	}

	template = linebot.NewCarouselTemplate(&columnList[0], &columnList[1], &columnList[2], &columnList[3], &columnList[4])


	//template = linebot.NewCarouselTemplate(
	//	linebot.NewCarouselColumn(
	//		"",
	//		"hoge",
	//		"fuga",
	//		linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
	//		//linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", "", ""),
	//	),
	//	linebot.NewCarouselColumn(
	//		"",
	//		"hoge",
	//		"fuga",
	//		linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
	//		linebot.NewMessageTemplateAction("Say message", "Rice=米"),
	//	),
	//)
	return template
}

func sendCarouselMessage(event *linebot.Event, template *linebot.CarouselTemplate){
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Carousel alt text", template)).Do(); err != nil {
		log.Println(err)
	}
}
