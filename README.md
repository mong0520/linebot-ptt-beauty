# 表特看看 - Line 聊天機器人 for PTT Beauty

## For User

### 掃描 QR Code 或點選連結
[<img src="resource/qr_code.png">](https://line.me/R/ti/p/SFXWQpzdaY)


---

## For Developer

### 本機測試

請先前往 [Line Message API](https://developers.line.biz/en/services/messaging-api/) 申請 AccessKey & Secret Key 以及相關設定，在此不贅述
```
# 1) 設定環境變數
copy .env.template .env
vim .env

# 2) Compile docker image
make build NAMESPAMCE=YOUR_DOCKER_NAMESPACE

# 3) 啟動 Linebot + MongDB
make dev

# 4) 使用 ngork 產生 https endpoint, 並至 Line Message API 後台設定 callback url 為 https://YOUR_NGROK_URL/callback
ngrok http 5000
```

### 資料注入
使用 https://github.com/mong0520/ptt-web-crawler/blob/master/run.sh，將資料注入 mongoDB 即可


### 佈署
可以佈署至 Heroku 測試使用，需要在 heroku dashboard 中設定 `ChannelAccessToken`, `ChannelSecret`, `MongoDBHostPort`, 與 `PORT` 參數，同 `.env.tempalte` 中之設定
```
heroku login
make push
make release
```

### 截圖

* 功能選單

<img src="resource/screen1.jpg" height="480">


* 熱門照片

<img src="resource/screen2.jpg" height="480">


* 對話直接搜尋

<img src="resource/screen3.jpg" height="480">

### 待辨清單:
