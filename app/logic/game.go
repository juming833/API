package logic

import (
	"IP/app/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetGame(c *gin.Context) {
	id := c.Query("id")
	Id, _ := strconv.Atoi(id)
	ret := model.GetGame(Id)
	c.JSON(200, gin.H{"games": ret})

}
