package model

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Rdb *redis.Client
var Conn *gorm.DB

func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "mysql_mn5dKj", "192.168.18.245:3306", "beego_admin")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库链接错误")
		panic(err)
	}
	// 设置数据表
	Conn = conn
}
func NewRdb() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.18.245:6379",
		Password: "jhkdjhkjdhsIUTYURTU_M4RAd5",
		DB:       0,
	})
	// 发送 Ping 命令来检查连接
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		// 连接失败，进行错误处理
		log.Printf("Redis连接失败: %s", err)
		return
	}

	Rdb = rdb
	log.Println("Redis连接成功")
	return
}
func Close() {
	_ = Rdb.Close()
	db, _ := Conn.DB()
	_ = db.Close()
}
