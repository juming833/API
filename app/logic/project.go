package logic

import (
	"IP/app/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetIPS(c *gin.Context) {
	ret := model.GetIPS()
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
