package app

import (
	"IP/app/model"
)

func Start() {
	model.NewMysql()
	model.NewRdb()
	//defer func() {
	//	model.Close()
	//}()
}
