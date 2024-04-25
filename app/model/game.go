package model

import "fmt"

func GetGame(id int) Game {
	var ret Game
	if err := Conn.Table("game").Where("id = ?", id).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}
