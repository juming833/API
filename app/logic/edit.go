package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Info struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
type Res struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"msg"`
}
type IDS struct {
	Id       []int64 `json:"id"`
	Password string  `json:"password"`
}

// EditPassword 修改密码
func EditPassword(c *gin.Context, key string) {
	//var id IDS
	//if err := c.ShouldBindJSON(&id); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//var usernames []string
	//for _, id := range id.Id {
	//	res := model.GetInfo(id)
	//	usernames = append(usernames, res.Username)
	//}
	//
	//c.JSON(http.StatusOK, gin.H{"usernames": usernames})
	username := c.PostForm("username")
	password := c.PostForm("password")
	formData := url.Values{}
	formData.Set("key", key)
	//for _, username := range usernames {
	//	formData.Set("username", username)
	//}
	formData.Set("username", username)
	formData.Set("password", password)
	Url := fmt.Sprintf("http://quzhous501.ipdog.cn:3999/admin/user/edit?key=%s", key)
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
