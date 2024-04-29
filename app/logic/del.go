package logic

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type DelRes struct {
	Code    int    `json:"code"`
	Data    int    `json:"data"` //指删除了多少条数据
	Message string `json:"msg"`
}

// DelSingle 删除指定用户
func DelSingle(c *gin.Context, key string) {
	username := c.Query("username")
	url := fmt.Sprintf("http://quzhous501.ipdog.cn:3999/admin/user/del?key=%s&username=%s", key, username)
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
	var responseData Res
	if err := json.Unmarshal(body, &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseData)
}

// DelAllUser 删除所有用户
func DelAllUser(c *gin.Context, key string) {
	url := fmt.Sprintf("http://quzhous501.ipdog.cn:3999/admin/user/delall?key=%s", key)
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
	var responseData DelRes
	if err := json.Unmarshal(body, &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseData)
}
