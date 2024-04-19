package logic

import (
	"IP/app/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
}

func GetYu(c *gin.Context) {
	var idlist IdList
	if err := c.ShouldBindJSON(&idlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//creatTimeUnix, _ := strconv.ParseInt(idlist.CreatTime, 10, 64)
	//CreatTime := time.Unix(creatTimeUnix, 0)
	//deleteTimeUnix, _ := strconv.ParseInt(idlist.DeleteTime, 10, 64)
	//DelTime := time.Unix(deleteTimeUnix, 0)
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
	jsonData, err := model.Rdb.Get(c, "quzhous501").Bytes()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var data map[string]map[string]string
	if err := json.Unmarshal(jsonData, &data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ret := model.GetAllIPs()
	finalIPs := make(map[string]string)
	for _, ip := range ret {
		if strings.HasPrefix(ip, "quzhous501") {
			numPart := strings.TrimPrefix(ip, "quzhous501"+"p")
			numPart = strings.Split(numPart, ".")[0]
			for jsonKey, jsonValue := range data {
				if strings.HasPrefix(jsonKey, "pppoe-out") {
					pppoeNum := strings.TrimPrefix(jsonKey, "pppoe-out")
					if pppoeNum == numPart {
						finalIPs[ip] = jsonValue["ip"]
						break
					}
				}
			}
		}
	}
	res := model.GetUPS()
	finalResponse := make([]map[string]string, 0)
	for _, project := range res {
		if _, ok := validIPIDs[project.Id]; ok {
			if ip, exists := finalIPs[project.Ip]; exists {
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
				ipWithPort := ip
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
	}
	var textData string
	for _, data := range finalResponse {
		ip := data["ip"]
		username := data["username"]
		password := data["password"]
		deleteTime := data["delete_time"]
		line := fmt.Sprintf("%s %s %s %s\n", ip, username, password, deleteTime)
		textData += line
	}
	var apiLists []model.ApiList
	for userID, ips := range userIPMap {
		combinedIPs := strings.Join(ips, ", ")
		creatTimeUnix, _ := strconv.ParseInt(creatTimes[0], 10, 64) // Assuming all entries have the same timestamps
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
	c.String(http.StatusOK, textData)
}
