package logic

import (
	"IP/app/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Info2Struct struct {
	Type string `json:"type"`
}
type Short struct {
	ProxyType int    `json:"proxy_type"`
	GameId    int    `json:"game_id"`
	MsId      int    `json:"ms_id"`
	UserId    int    `json:"user_id"`
	Count     int    `json:"buy_count"`
	ModeId    int    `json:"mode_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
	Url  string `json:"url"`
	Wait int    `json:"wait"`
}

// Order 购买节点
func Order(c *gin.Context, key string) {
	var request Short
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	originalCount := request.Count
	var info model.Game
	info = model.GetGame(request.GameId)
	if info.Info2 != "" {
		var info2 Info2Struct
		err := json.Unmarshal([]byte(info.Info2), &info2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error: " + err.Error()})
			return
		}
		if info2.Type == "short" { //判断类型，走短效接口
			var network []model.IpList
			network = model.GetIPS(request.MsId)
			var networkIPs []string
			for _, ipList := range network {
				networkIPs = append(networkIPs, ipList.Network)
			}
			//request.Count = 1
			formData := url.Values{}
			formData.Set("key", key)
			formData.Set("enable", fmt.Sprintf("%v", true))                                    //状态
			formData.Set("expiration", time.Now().Add(duration).Format("2006-01-02 15:04:05")) //过期时间
			formData.Set("info", fmt.Sprintf("%d", request.UserId))                            //用户id
			formData.Set("maxconn", fmt.Sprintf("%d", info.MaxConn))                           //最大连接数
			formData.Set("network", strings.Join(networkIPs, ","))
			formData.Set("password", request.Password)
			formData.Set("recv", fmt.Sprintf("%d", info.DwData)) //下行宽带
			formData.Set("send", fmt.Sprintf("%d", info.UpData)) //上行宽带
			formData.Set("username", request.Username)
			formData.Set("usertype", fmt.Sprintf("%d", 1))
			Url := fmt.Sprintf("http://quzhous501.ipdog.cn:3999/admin/user/add?key=%s", key)
			executeRequest1(c, formData, Url)
			//request.MsId = 68
			Url2 := "http://s5test.jianyunip.com/api/account/buy"
			executeRequest2(c, request, Url2)
			res := model.ShortUserProject{
				//UserId:   request.UserId,
				BuyCount: originalCount,
				//MsId:       request.MsId,
				DeleteTime: time.Now().Add(duration).Format("2006-01-02 15:04:05"),
				CreatTime:  time.Now().Format("2006-01-02 15:04:05"),
			}
			if err := model.AddShort(res); err != nil {
				c.JSON(http.StatusOK, gin.H{"error": "写入数据库错误！"})
				return
			}
		} else { //长效购买
			Url := "http://s5test.jianyunip.com/api/account/buy"
			executeRequest2(c, request, Url)
		}
	} else {
		Url := "http://s5test.jianyunip.com/api/account/buy"
		executeRequest2(c, request, Url)
	}
}

// 表单请求
func executeRequest1(c *gin.Context, formData url.Values, url string) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HTTP request: " + err.Error()})
		return
	}
	token := c.GetHeader("Authorization")
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
}

// json请求
func executeRequest2(c *gin.Context, data interface{}, url string) {
	jsonData, err := json.Marshal(data)
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
	c.JSON(http.StatusOK, gin.H{"message2": responseData})
}
