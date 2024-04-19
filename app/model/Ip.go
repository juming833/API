package model

import (
	"fmt"
	"gorm.io/gorm"
)

func GetIPS() []IpList {
	ret := make([]IpList, 0)
	if err := Conn.Table("ip_list").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

type Result struct {
	IpList
	ManServerName string
	Addr          string
	SocksPort     int
	HttpPort      int
	SsuPort       int
}
type Res struct {
	UserProject
	Ip        string
	Id        int
	SocksPort int
	HttpPort  int
	SsuPort   int
}

func GetIP() []Result {
	var results []Result
	if err := Conn.Table("ip_list").
		Select("ip_list.*, man_server.name as man_server_name, man_server.addr as addr,man_server.socks_port as socks_port,man_server.http_port as http_port,man_server.ssu_port as ssu_port").
		Joins("join man_server on ip_list.ms_i_d_id = man_server.id").
		Scan(&results).Error; err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return results
}

func GetAllIPs() []string {
	var ips []string
	if err := Conn.Table("ip_list").Pluck("ip", &ips).Error; err != nil {
		fmt.Printf("err: %s", err.Error())
	}
	return ips
}
func GetUP() []UserProject {
	ret := make([]UserProject, 0)
	if err := Conn.Table("user_project").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}
func GetUPS() []Res {
	var results []Res
	subQuery := Conn.Table("ip_list").Select("ms_i_d_id").Where("ip_list.id = user_project.ip_id")
	if err := Conn.Table("user_project").
		Select("user_project.*, ip_list.ip as ip, man_server.socks_port as socks_port, man_server.http_port as http_port, man_server.ssu_port as ssu_port").
		Joins("JOIN ip_list ON user_project.ip_id = ip_list.id").
		Joins("JOIN man_server ON man_server.id = (?)", subQuery).
		Scan(&results).Error; err != nil {
	}
	return results
}
func AddApiList(Info ApiList) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Info).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
