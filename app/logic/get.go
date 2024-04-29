package logic

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type IpList struct {
	Code int    `json:"code"`
	Data Data   `json:"data"`
	Msg  string `json:"msg"`
}

type Data struct {
	List []Item `json:"list"`
}

type Item struct {
	Network    string `json:"network"`
	OuternetIP string `json:"outernetip"`
}

// GetIpList 获取IP列表
func GetIpList(c *gin.Context, key string) {
	url := fmt.Sprintf("http://quzhous501.ipdog.cn:3999/admin/getip/list?key=%s&update=%v", key, true)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body: " + err.Error()})
		return
	}
	var responseData IpList
	if err := json.Unmarshal(body, &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseData)
}
