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
	"strconv"
	"strings"
	"time"
)

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
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
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
	formData.Set("enable", fmt.Sprintf("%v", true))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败: " + err.Error()})
		return
	}
	token := c.GetHeader("Authorization")
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求失败: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取失败: " + err.Error()})
		return
	}
	var responseData Response
	if err := json.Unmarshal(body, &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化错误: " + err.Error()})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"message": responseData})
	id := GetUID()
	if responseData.Code != 200 {
		c.JSON(0, gin.H{"error1": responseData.Msg})
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
	//filePath := strconv.Itoa(int(id)) + ".txt"
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
	dataKey := "userData:" + strconv.FormatInt(id, 10)
	separator := c.PostForm("customValue")

	formattedData := formatData(Format, networks, Port, Username, Password, Expiration, separator)
	err = model.Rdb.Set(c, dataKey, formattedData, -1).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis 持久化错误：" + err.Error()})
		return
	}
	idStr := strconv.Itoa(int(id))
	UrlStr := "http://192.168.18.142:8081/api/get?id=" + idStr + "&format=" + format + "&customValue=" + separator
	c.JSON(http.StatusOK, gin.H{"url": UrlStr})
}

func GetData(c *gin.Context) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	id := c.Query("id")
	format := c.Query("format")
	separator := c.Query("customValue")
	Format, _ := strconv.Atoi(format)
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "参数错误！"})
		return
	}
	dataKey := "userData:" + id
	data, err := model.Rdb.Get(c, dataKey).Bytes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis检索错误: " + err.Error()})
		return
	}
	var (
		networks   []string
		port       string
		username   string
		password   string
		expiration string
	)
	switch Format {
	case 1:
		var dataMap map[string]interface{}
		if err := json.Unmarshal(data, &dataMap); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析 JSON 错误: " + err.Error()})
			return
		}
		if ipList, ok := dataMap["IPs"].([]interface{}); ok {
			for _, ip := range ipList {
				if ipStr, ok := ip.(string); ok {
					networks = append(networks, ipStr)
				}
			}
		}
		port, _ = dataMap["Port"].(string)
		username, _ = dataMap["Username"].(string)
		password, _ = dataMap["Password"].(string)
		expiration, _ = dataMap["ExTime"].(string)
	case 2:
		dataStr := string(data)
		lines := strings.Split(dataStr, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.Split(line, " ")
			if len(parts) > 0 {
				networks = append(networks, parts[0])
				port = parts[1]
				username = parts[2]
				password = parts[3]
				expiration = parts[4]
			}
		}

	case 3:
		dataStr := string(data)
		lines := strings.Split(dataStr, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) > 0 {
				networks = append(networks, parts[0])
				port = parts[1]
				username = parts[2]
				password = parts[3]
				expiration = parts[4]
			}
		}
	case 4:
		dataStr := string(data)
		lines := strings.Split(dataStr, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.Split(line, "+")
			if len(parts) > 0 {
				networks = append(networks, parts[0])
				port = parts[1]
				username = parts[2]
				password = parts[3]
				expiration = parts[4]
			}
		}
	}
	keys := GetRedisKeys(c)
	Data := make(map[string]map[string]map[string]string)
	for _, key := range keys {
		jsonData, err := model.Rdb.Get(c, key).Bytes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var tempData map[string]map[string]string
		if err := json.Unmarshal(jsonData, &tempData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
			return
		}
		Data[key] = tempData
	}

	finalIPs := make(map[string]string)
	for _, ip := range networks {
		parts := strings.Split(ip, ".")
		lastPart := parts[len(parts)-1]
		for _, ipInfoMap := range Data {
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
	if format == "1" {
		for _, network := range networks {
			if newIp, exists := finalIPs[network]; exists {
				updatedNetworks = append(updatedNetworks, fmt.Sprintf("%s %s %s %s %s", newIp, port, username, password, expiration))
			} else {
				updatedNetworks = append(updatedNetworks, fmt.Sprintf("%s %s %s %s %s", network, port, username, password, expiration))
			}
		}
	} else if format == "2" {
		for _, network := range networks {
			if newIp, exists := finalIPs[network]; exists {
				updatedNetworks = append(updatedNetworks, fmt.Sprintf("%s:%s %s %s %s", newIp, port, username, password, expiration))
			} else {
				updatedNetworks = append(updatedNetworks, fmt.Sprintf("%s:%s %s %s %s", network, port, username, password, expiration))
			}
		}
	} else if format == "3" {
		for _, network := range networks {
			if newIp, exists := finalIPs[network]; exists {
				updatedNetworks = append(updatedNetworks, fmt.Sprintf("%s|%s|%s|%s|%s", newIp, port, username, password, expiration))
			} else {
				updatedNetworks = append(updatedNetworks, fmt.Sprintf("%s|%s|%s|%s|%s", network, port, username, password, expiration))
			}
		}
	} else if format == "4" {
		for _, network := range networks {
			if newIp, exists := finalIPs[network]; exists {
				updatedNetworks = append(updatedNetworks, fmt.Sprintf("%s%s%s%s%s%s%s%s%s", newIp, separator, port, separator, username, separator, password, separator, expiration))
			} else {
				updatedNetworks = append(updatedNetworks, fmt.Sprintf("%s%s%s%s%s%s%s%s%s", network, separator, port, separator, username, separator, password, separator, expiration))
			}
		}
	} else {
		return
	}
	if format == "1" {
		responseData, _ := json.Marshal(map[string]interface{}{
			"ips":        updatedNetworks,
			"port":       port,
			"username":   username,
			"password":   password,
			"expiration": expiration,
		})
		c.Data(http.StatusOK, "application/json", responseData)
	} else {
		responseText := strings.Join(updatedNetworks, "\n")
		c.Data(http.StatusOK, "text/plain", []byte(responseText))
	}
}

func formatData(format int, networks []string, port, username, password, expiration, customSeparator string) string {
	var fileData string
	switch format {
	case 1:
		jsonData, _ := json.Marshal(map[string]interface{}{"IPs": networks, "Port": port, "Username": username, "Password": password, "ExTime": expiration})
		fileData = string(jsonData)
	case 2:
		for _, ip := range networks {
			fileData += fmt.Sprintf("%s %s %s %s %s\n", ip, port, username, password, expiration)
		}
	case 3:
		for _, ip := range networks {
			fileData += fmt.Sprintf("%s|%s|%s|%s|%s\n", ip, port, username, password, expiration)
		}
	case 4:
		for _, ip := range networks {
			fileData += fmt.Sprintf("%s%s%s%s%s%s%s%s%s\n", ip, "+", port, "+", username, "+", password, "+", expiration)
		}
		//fileData += customSeparator + "\n"
	}
	return fileData
}
