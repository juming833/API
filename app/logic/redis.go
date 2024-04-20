package logic

import (
	"IP/app/model"
	"github.com/gin-gonic/gin"
)

func GetRedisKeys(c *gin.Context) []string {
	keys := model.Rdb.Keys(c, "*").Val()
	//c.JSON(http.StatusOK, gin.H{
	//	"keys": keys,
	//})
	return keys
}
