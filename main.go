package main

import (
	"IP/app"
	"IP/app/logic"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	app.Start()
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	key := "jianyun"
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000"}
	//r.Use(cors.New(config))
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	r.POST("/api/newIpUpdate", logic.Write) //拨号
	r.GET("/getIPS", logic.GetIPS)
	r.GET("/get", logic.Read)
	r.GET("/get/city", logic.GetCity)
	r.GET("/getsss", logic.Getsss)
	r.GET("/get/citys", logic.GetCitys)
	r.GET("/get/Pr", logic.GetPr)

	//r.POST("/getNode", logic.Node)
	r.POST("/getYu", logic.GetYu)
	//r.POST("/getUrl", logic.GetUrl)
	r.POST("/short", func(c *gin.Context) {
		logic.Order(c, key)
	})
	r.POST("/getUrl", func(c *gin.Context) {
		logic.GetUrl(c, key)
	})
	r.GET("/getUP", logic.GetUPS)
	r.GET("/getIP", logic.GetIP)
	r.GET("/Api/Data", logic.GetTextData)
	r.GET("/Allip", logic.GetAllIP)
	r.GET("/Game", logic.GetGame)
	r.GET("/getIps", logic.GetIps)
	r.GET("/Info", logic.GetInfo)
	//r.GET("/redis", logic.GetRedisKeys)
	r.POST("/user/edit", func(c *gin.Context) {
		logic.EditPassword(c, key)
	})
	r.GET("/user/del", func(c *gin.Context) {
		logic.DelSingle(c, key)
	})
	r.GET("/user/delall", func(c *gin.Context) {
		logic.DelAllUser(c, key)
	})
	r.GET("/getip/list", func(c *gin.Context) {
		logic.GetIpList(c, key)
	})
	r.POST("/black/add", func(c *gin.Context) {
		logic.AddBlack(c, key)
	})
	r.GET("/black/delall", func(c *gin.Context) {
		logic.DelAllBlack(c, key)
	})
	r.GET("/api/get", logic.GetData)
	if err := r.Run(":8081"); err != nil {

		panic("gin 启动失败")
	}

}
