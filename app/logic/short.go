package logic

import (
	"IP/app/model"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

type Short struct {
	ProxyType int `json:"proxy_type"`
	GameId    int `json:"game_id"`
	MsId      int `json:"ms_id"`
	UserId    int `json:"user_id"`
	Count     int `json:"buy_count"`
	ModeId    int `json:"mode_id"`
}
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
	Url  string `json:"url"`
	Wait int    `json:"wait"`
}

func Order(c *gin.Context) {
	var request Short
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//前端传的数量
	originalCount := request.Count
	request.Count = 1
	url := "http://s5test.jianyunip.com/api/account/buy"
	jsonData, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling JSON: " + err.Error()})
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HTTP request: " + err.Error()})
		return
	}
	token := c.GetHeader("Authorization")
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MjY3NDMsInVzZXJuYW1lIjoiNjY2NjY2IiwibW9iaWxlIjoiMTExMTExMTExMTEiLCJwYXNzd29yZCI6ImYzNzllYWYzYzgzMWIwNGRlMTUzNDY5ZDFiZWMzNDVlIiwic3RhdHVzIjoxLCJkZXNjcmlwdGlvbiI6ImRhaWxhbyIsIndhbGxldCI6NDk5OTAsImNyZWF0ZV90aW1lIjoxNzE0MDMwNTMwLCJwb3BfdXNlciI6MCwiVXNlckxldmVsIjp7ImlkIjoyLCJuYW1lIjoi55m96ZO25Lya5ZGYIiwiZGVzY3JpcHRpb24iOiLljZXmrKHlhYXlgLzotoXov4c1MDAw5YWD5LqrOS445oqY5LyY5oOgIiwiYW1vdW50Ijo1MDAwLCJkaXNjb3VudCI6MC45OCwic3RhdHVzIjoxLCJjcmVhdGVfdGltZSI6MTU4Nzg3OTg3MSwidXBkYXRlX3RpbWUiOjE2NTcwOTQxNDcsImRlbGV0ZV90aW1lIjowfSwiUHJvRmlsZSI6eyJpZCI6NDksImxhc3RfbmFtZSI6IiIsImlkX2NhcmQiOiIiLCJjZXJ0aWZ5X2lkIjoiIiwiY2VydGlmeV9ubyI6IiIsInN0YXRlIjowfSwiZXhwIjoxNzE0MjAzNDA0LCJpYXQiOjE3MTQwMzA2MDQsImlzcyI6IjY2NjY2NiJ9.4gai5N9IJaPuUxcGAH2DIVq8KzqCMF4kb9isVeIEY2Q")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make HTTP request: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body: " + err.Error()})
		return
	}
	var responseData Response
	if err := json.Unmarshal(body, &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": responseData})

	var duration time.Duration
	switch request.ModeId {
	case 1:
		duration = 24 * time.Hour
	case 2:
		duration = 7 * 24 * time.Hour
	case 3:
		duration = 30 * 24 * time.Hour
	case 4:
		duration = 90 * 24 * time.Hour
	case 5:
		duration = 365 * 24 * time.Hour
	}
	//UserId, _ := strconv.Atoi(request.UserId)
	info := model.ShortUserProject{
		Type:       request.ProxyType,
		GameId:     request.GameId,
		UserId:     request.UserId,
		BuyCount:   originalCount,
		ModeId:     request.ModeId,
		MsId:       request.MsId,
		DeleteTime: time.Now().Add(duration).Format("2006-01-02 15:04:05"),
		CreatTime:  time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := model.AddShort(info); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "写入数据库错误！"})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"message": "成功处理请求"})
}
