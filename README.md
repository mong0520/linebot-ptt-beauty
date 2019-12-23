# 表特看看 - Line 聊天機器人 for PTT Beauty

## For User

### 掃描 QR Code 或點選連結
[<img src="resource/qr_code.png">](https://line.me/R/ti/p/SFXWQpzdaY)


---

## For Developer

### 開發環境

```
go get github.com/mong0520/linebot-ptt
govendor sync
```


### 本機測試

```
# 1) 啟動 MongDB
Follow the instruction in db/

# 2) 啟動 Linebot
export PORT=${PORT}
export ChannelSecret=${ChannelSecret}
export ChannelAccessToken=${ChannelAccessToken}

go run main.go

# 3) 設定 Https 轉發
ngrok http 8080

# 4) 設定 Linebot webhook

```


### 截圖

* 功能選單

<img src="resource/screen1.jpg" height="480">


* 熱門照片

<img src="resource/screen2.jpg" height="480">


* 對話直接搜尋

<img src="resource/screen3.jpg" height="480">

### 待辨清單:
