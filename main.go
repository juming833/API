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
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))
	r.POST("/api/newIpUpdate", logic.Write)
	r.GET("/getIPS", logic.GetIPS)
	r.GET("/get", logic.Read)
	//r.POST("/getNode", logic.Node)
	r.POST("/getYu", logic.GetYu)
	r.GET("/getUP", logic.GetUPS)
	r.GET("/getIP", logic.GetIP)
	r.GET("/Api/Data", logic.GetTextData)
	r.GET("/Allip", logic.GetAllIP)
	//r.GET("/redis", logic.GetRedisKeys)

	if err := r.Run(":8081"); err != nil {
		panic("gin 启动失败")
	}

}
