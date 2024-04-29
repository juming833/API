package logic

import (
	"IP/app/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetIPS(c *gin.Context) {
	id := c.Query("id")
	Id, _ := strconv.Atoi(id)
	ret := model.GetIPS(Id)
	c.JSON(http.StatusOK, gin.H{"ip": ret})
}
func GetUPS(c *gin.Context) {
	ret := model.GetUPS()
	c.JSON(http.StatusOK, gin.H{"project": ret})
}
func GetIP(c *gin.Context) {
	ret := model.GetIP()
	c.JSON(http.StatusOK, gin.H{"ip": ret})
}
func GetAllIP(c *gin.Context) {
	ret := model.GetAllIPs()
	c.JSON(http.StatusOK, gin.H{"ip": ret})
}
func GetGame(c *gin.Context) {
	id := c.Query("id")
	Id, _ := strconv.Atoi(id)
	ret := model.GetGame(Id)
	c.JSON(200, gin.H{"games": ret})

}
func GetRedisKeys(c *gin.Context) []string {
	keys := model.Rdb.Keys(c, "*").Val()
	return keys
}
func GetInfo(c *gin.Context) {
	Id := c.Query("id")
	id, _ := strconv.Atoi(Id)
	ret := model.GetInfo(int64(id))
	c.JSON(http.StatusOK, gin.H{"ip": ret})
}
func GetCity(c *gin.Context) {
	ret := model.GetCity()
	c.JSON(http.StatusOK, gin.H{"city": ret})
}
func GetCitys(c *gin.Context) {
	A := model.GetCitys()
	B := model.GetPr()
	ret := model.MergeData(B, A)
	fmt.Println(ret)
	c.JSON(http.StatusOK, gin.H{"citys": ret})
}
func GetPr(c *gin.Context) {
	ret := model.GetPr()
	fmt.Println(ret)
	c.JSON(http.StatusOK, gin.H{"Pr": ret})
}
func GetIps(c *gin.Context) {
	Id := c.Query("id")
	id, _ := strconv.Atoi(Id)
	ret := model.GetIps(id)
	fmt.Println(ret)
	c.JSON(http.StatusOK, gin.H{"Pr": ret})
}
func Getsss(c *gin.Context) {
	Id := c.Query("id")

	ret := model.Getsss(Id)
	fmt.Println(ret)
	c.JSON(http.StatusOK, gin.H{"Pr": ret})
}
