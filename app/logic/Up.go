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
	"os"
	"strconv"
	"strings"
	"time"
)

type DaTa struct {
	Ips      []string `json:"ips"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

func GetUrl(c *gin.Context, key string) {
	MsID := c.PostForm("cityId")
	port := c.PostForm("proxyType")
	portInt, _ := strconv.Atoi(port)
	format := c.PostForm("geishi")
	Format, _ := strconv.Atoi(format)
	Username := c.PostForm("username")
	Password := c.PostForm("password")
	Expiration := c.PostForm("expiration")
	UserId := c.PostForm("userid")
	Num := c.PostForm("num")
	num, _ := strconv.Atoi(Num)
	var network []model.IpList
	network = model.Getsss(MsID)
	var networks []string
	for _, ipList := range network {
		networks = append(networks, ipList.Network)
	}
	keys := GetRedisKeys(c)
	data := make(map[string]map[string]map[string]string)
	for _, key := range keys {
		jsonData, err := model.Rdb.Get(c, key).Bytes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var tempData map[string]map[string]string
		if err := json.Unmarshal(jsonData, &tempData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data[key] = tempData
	}
	finalIPs := make(map[string]string)
	for _, ip := range networks {
		parts := strings.Split(ip, ".")
		lastPart := parts[len(parts)-1]
		for _, ipInfoMap := range data {
			for jsonKeyInner, jsonValue := range ipInfoMap {
				if strings.HasPrefix(jsonKeyInner, "pppoe-out") {
					pppoeNum := strings.TrimPrefix(jsonKeyInner, "pppoe-out")
					if pppoeNum == lastPart {
						ipVal, ok := jsonValue["ip"]
						if ok {
							finalIPs[ip] = ipVal
							break
						}
					}
				}
			}
		}
	}
	updatedNetworks := make([]string, 0, len(networks))
	for _, network := range networks {
		if newIp, exists := finalIPs[network]; exists {
			updatedNetworks = append(updatedNetworks, newIp)
		} else {
			updatedNetworks = append(updatedNetworks, network)
		}
	}
	formData := url.Values{}
	formData.Set("key", key)
	formData.Set("enable", fmt.Sprintf("%v", true)) //状态
	//formData.Set("expiration", time.Now().Add(24time.Hour).Format("2006-01-02 15:04:05")) //过期时间
	formData.Set("expiration", Expiration) //过期时间
	formData.Set("info", UserId)           //用户id
	formData.Set("maxconn", "100")         //最大连接数
	formData.Set("network", strings.Join(networks, ","))
	formData.Set("password", Password)
	formData.Set("recv", "10000") //下行宽带
	formData.Set("send", "10000") //上行宽带
	formData.Set("username", Username)
	formData.Set("usertype", fmt.Sprintf("%d", 1))
	Url := fmt.Sprintf("http://quzhous501.ipdog.cn:3999/admin/user/add?key=%s", key)
	//executeRequest(c, formData, Url)
	req, err := http.NewRequest("POST", Url, bytes.NewBufferString(formData.Encode()))
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
	//c.JSON(http.StatusOK, gin.H{"message": responseData})
	id := GetUID()
	if responseData.Code != 200 {
		c.JSON(0, gin.H{"error": responseData.Msg})
		return
	}
	res := model.ShortUserProject{
		UserId:     UserId,
		Username:   Username,
		Password:   Password,
		BuyCount:   num,
		MsId:       MsID,
		ApiId:      strconv.FormatInt(id, 10),
		DeleteTime: Expiration,
		CreatTime:  time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := model.AddShort(res); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "写入数据库错误！"})
		return
	}
	//res = append(res, Username, Password)
	filePath := strconv.Itoa(int(id)) + ".txt"
	//data := fmt.Sprintln(res)
	type Data struct {
		Ip         string
		Username   string
		Password   string
		DeleteTime string
		Port       string
	}
	Port := ""
	switch portInt {
	case 1:
		Port = strconv.Itoa(3999)
	case 2:
		Port = strconv.Itoa(3998)
	}
	var dataList []Data
	for _, ip := range networks {
		data := Data{
			Port:       Port,
			Ip:         ip,
			Username:   Username,
			Password:   Password,
			DeleteTime: Expiration,
		}
		dataList = append(dataList, data)
	}
	var fileData string
	if Format == 1 {
		jsonData, _ := json.Marshal(map[string]interface{}{"IPs": networks, "Port": Port, "Username": Username, "Password": Password, "DeleteTime": Expiration})
		fileData = string(jsonData)
	} else {
		for _, ip := range updatedNetworks {
			switch Format {
			case 2:
				fileData += fmt.Sprintf("%s:%s %s %s %s\n", ip, Port, Username, Password, Expiration)
			case 3:
				fileData += fmt.Sprintf("%s|%s|%s|%s|%s\n", ip, Port, Username, Password, Expiration)
			case 4:
				separator := c.PostForm("customValue")
				fileData += fmt.Sprintf("%s%s%s%s%s%s%s%s%s\n", ip, separator, Port, separator, Username, separator, Password, separator, Expiration)
			}
		}
	}

	_ = os.WriteFile(filePath, []byte(fileData), 0644)
	expirationDuration := 300 * time.Second
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法写入文件"})
		return
	}
	time.AfterFunc(expirationDuration, func() { //设置文件自动过期清除
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("过期文件清除失败:", filePath)
		} else {
			fmt.Println("过期文件自动清除:", filePath)
		}
	})
	idStr := strconv.Itoa(int(id))
	UrlStr := "192.168.18.142:8081/Api/Data?id=" + idStr
	c.JSON(http.StatusOK, gin.H{"url": UrlStr,
		"code": 0,
	})
}
