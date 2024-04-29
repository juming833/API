package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// DelAllBlack 删除所有黑名单数据
func DelAllBlack(c *gin.Context, key string) {
	Url := fmt.Sprintf("http://quzhous501.ipdog.cn:3999/admin/black/delall?key=%s", key)
	response, err := http.Get(Url)
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

// AddBlack 添加黑名单
func AddBlack(c *gin.Context, key string) {
	data := c.PostForm("data")
	dataSlice := strings.Split(data, ",")
	formData := url.Values{}
	formData.Set("key", key)
	formData.Set("data", strings.Join(dataSlice, "\r\n"))
	Url := fmt.Sprintf("http://quzhous501.ipdog.cn:3999/admin/black/add?key=%s", key)
	req, err := http.NewRequest("POST", Url, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HTTP request: " + err.Error()})
		return
	}
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
	var responseData Res
	if err := json.Unmarshal(body, &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseData)
}
