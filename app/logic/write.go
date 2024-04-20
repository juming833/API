package logic

import (
	"IP/app/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type IPData struct {
	Ip        string `json:"ip"`
	Timestamp string `json:"timestamp"`
}

func Write(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	id := c.PostForm("id")
	data := c.PostForm("data")

	//IDList = append(IDList, id)
	//fmt.Println(IDList)
	if id == "" || data == "" {
		c.JSON(400, gin.H{"error": "接收数据的格式错误"})
		fmt.Println(1)
		return
	}
	exists, err := model.Rdb.Exists(c, id).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "查找redis失败"})
		return
	}
	if exists != 0 {
		c.JSON(400, gin.H{"error": "域名已经存在"})
		fmt.Println(2)
		return
	}
	// 分割data中的每个键值对
	keyValuePairs := strings.Split(data, "#")
	ipMap := make(map[string]IPData)

	for _, pair := range keyValuePairs {
		pair = strings.TrimSpace(pair)
		if pair != "" {
			parts := strings.Split(pair, "=")
			if len(parts) > 1 {
				key := parts[0]
				ip := parts[1]
				ipData := IPData{
					Ip:        ip,
					Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				}
				ipMap[key] = ipData
			}
		}
	}
	jsonData, err := json.Marshal(ipMap)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error marshalling JSON"})
		fmt.Println(err)
		return
	}
	if err := model.Rdb.Set(c, id, jsonData, 0).Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{"message": "存入redis成功"})
}

func Read(c *gin.Context) {
	//key := c.Query("key")
	jsonData, err := model.Rdb.Get(c, "quzhous501").Bytes()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	//var ipMap map[string]IPData
	//err = json.Unmarshal(jsonData, &ipMap)
	//if err != nil {
	//	c.JSON(500, gin.H{"error": "Error unmarshalling JSON"})
	//	fmt.Println(err)
	//	return
	//}
	//ipData, ok := ipMap[key]
	//if !ok {
	//	c.JSON(404, gin.H{"error": "指定的key不存在"})
	//	return
	//}
	//
	//c.JSON(200, gin.H{"ip": ipData.Ip, "timestamp": ipData.Timestamp})
	c.Data(200, "application/json", jsonData)
}

// Writer 存哈希表
//func Writer(c *gin.Context) {
//	if err := c.Request.ParseForm(); err != nil {
//		c.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//	id := c.PostForm("id")
//	data := c.PostForm("data")
//	IDs := id + "s"
//	if id == "" || data == "" {
//		c.JSON(400, gin.H{"error": "接收数据的格式错误"})
//
//		return
//	}
//	exists, err := model.Rdb.Exists(c, id).Result()
//	if err != nil {
//		c.JSON(500, gin.H{"error": "查找redis失败"})
//		return
//	}
//	if exists != 0 {
//		c.JSON(400, gin.H{"error": "域名已经存在"})
//		return
//	}
//	keyValuePairs := strings.Split(data, "#")
//	for _, pair := range keyValuePairs {
//		pair = strings.TrimSpace(pair)
//		if pair != "" {
//			parts := strings.Split(pair, "=")
//			if len(parts) > 1 {
//				key := parts[0]
//				ip := parts[1]
//				if err := model.Rdb.HSet(c, IDs, key, ip).Err(); err != nil {
//					fmt.Printf("Failed to set hash key %s in hash %s: %s\n", key, id, err)
//					continue
//				}
//			}
//		}
//	}
//	c.JSON(200, gin.H{"message": "存入redis成功"})
//}
