package logic

import (
	"IP/app/model"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type IdList struct {
	IdList     string `json:"id"`
	UserID     string `json:"user_id"`
	CreatTime  string `json:"creat_time"`
	DeleteTime string `json:"delete_time"`
	IPID       string `json:"ip_id"`
	FormatID   int    `json:"format_id"`
}

func GetYu(c *gin.Context) {
	var idlist IdList
	if err := c.ShouldBindJSON(&idlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ids := strings.Split(idlist.IdList, ",")
	userIDs := strings.Split(idlist.UserID, ",")
	creatTimes := strings.Split(idlist.CreatTime, ",")
	deleteTimes := strings.Split(idlist.DeleteTime, ",")
	ipIDs := strings.Split(idlist.IPID, ",")
	if len(ids) != len(userIDs) || len(userIDs) != len(creatTimes) || len(creatTimes) != len(deleteTimes) || len(deleteTimes) != len(ipIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误，数组长度不匹配"})
		return
	}
	userIPMap := make(map[string][]string)
	for i, userID := range userIDs {
		userIPMap[userID] = append(userIPMap[userID], ipIDs[i])
	}
	validIPIDs := make(map[int]bool)
	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid id value: %s", idStr)})
			return
		}
		validIPIDs[id] = true
	}
	keys := GetRedisKeys(c)
	data := make(map[string]map[string]map[string]string)

	// 遍历每个键并获取对应的数据
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

	ret := model.GetAllIPs()
	finalIPs := make(map[string]string)
	for _, ip := range ret {
		for prefix, jsonValueMap := range data {
			if strings.HasPrefix(ip, prefix) {
				numPart := strings.TrimPrefix(ip, prefix+"p")
				numPart = strings.Split(numPart, ".")[0]
				for jsonKeyInner, jsonValue := range jsonValueMap {
					if strings.HasPrefix(jsonKeyInner, "pppoe-out") {
						pppoeNum := strings.TrimPrefix(jsonKeyInner, "pppoe-out")
						if pppoeNum == numPart {
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
	}

	res := model.GetUPS()
	finalResponse := make([]map[string]string, 0)
	for _, project := range res {
		IP := project.Ip
		if resolvedIP, exists := finalIPs[project.Ip]; exists {
			IP = resolvedIP
		}
		if _, ok := validIPIDs[project.Id]; ok {
			if _, ok := validIPIDs[project.Id]; ok {
				if ip, exists := finalIPs[project.Ip]; exists {
					IP = ip
				}
			}

			deleteTime := strconv.Itoa(project.DeleteTime)
			tm, _ := strconv.ParseInt(deleteTime, 10, 64)
			port := ""
			switch project.Type {
			case 1:
				port = strconv.Itoa(project.SocksPort)
			case 2:
				port = strconv.Itoa(project.HttpPort)
			case 3:
				port = strconv.Itoa(project.SsuPort)
			}
			ipWithPort := IP
			if port != "" {
				ipWithPort += ":" + port
			}
			finalResponse = append(finalResponse, map[string]string{
				"ip":          ipWithPort,
				"username":    project.Username,
				"password":    project.Password,
				"delete_time": time.Unix(tm, 0).Format("2006-01-02 15:04:05"),
			})
		}
	}

	var textData string
	for _, data := range finalResponse {
		ip := data["ip"]
		username := data["username"]
		password := data["password"]
		deleteTime := data["delete_time"]
		if idlist.FormatID == 1 {
			line := fmt.Sprintf("%s %s %s %s\n", ip, username, password, deleteTime)
			textData += line
		} else if idlist.FormatID == 2 {
			line := fmt.Sprintf("%s| %s| %s| %s\n", ip, username, password, deleteTime)
			textData += line
		}
	}
	var apiLists []model.ApiList
	for userID, ips := range userIPMap {
		combinedIPs := strings.Join(ips, ",")
		creatTimeUnix, _ := strconv.ParseInt(creatTimes[0], 10, 64)
		deleteTimeUnix, _ := strconv.ParseInt(deleteTimes[0], 10, 64)
		UserIDs, _ := strconv.Atoi(userID)
		api := model.ApiList{
			UserId:     UserIDs,
			IpIdList:   combinedIPs,
			CreatTime:  time.Unix(creatTimeUnix, 0),
			DeleteTime: time.Unix(deleteTimeUnix, 0),
		}
		apiLists = append(apiLists, api)
	}
	for _, api := range apiLists {
		if err := model.AddApiList(api); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("写入数据库错误：%v", err)})
			return
		}
	}
	//c.String(http.StatusOK, textData)
	id := GetUID()
	filePath := strconv.Itoa(int(id)) + ".txt"
	err := os.WriteFile(filePath, []byte(textData), 0644)
	expirationDuration := 300 * time.Second
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法写入文件"})
		return
	}
	time.AfterFunc(expirationDuration, func() {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("过期文件清除失败:", filePath)
		} else {
			fmt.Println("过期文件自动清除:", filePath)
		}
	})
	idStr := strconv.Itoa(int(id))
	c.JSON(http.StatusOK, gin.H{"id": idStr})

}
func GetTextData(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "参数错误！"})
		return
	}
	filePath := id + ".txt"
	data, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "id不存在或者已过期！"})
		return
	}
	c.String(http.StatusOK, string(data))
}

var snowNode *snowflake.Node

func GetUID() int64 {
	if snowNode == nil {
		snowNode, _ = snowflake.NewNode(1)
	}
	return snowNode.Generate().Int64()
}
