package model

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func GetIPS(id int) []IpList {
	var ret []IpList
	if err := Conn.Table("ip_list").Where("ms_i_d_id = ?", id).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetIps(id int) []string {
	var ret []IpList
	if err := Conn.Table("ip_list").Select("ip").Where("ms_i_d_id = ?", id).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	ips := make([]string, len(ret))
	for i, item := range ret {
		ips[i] = item.Ip
	}
	return ips
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
	msId      int
}
type CITY struct {
	MsId   int    `json:"id"`
	MsName string `json:"city"`
	//CityId   int    `json:"CityId"`
	CityName string `json:"province"`
}

func GetIP() []Result {
	var results []Result
	if err := Conn.Table("ip_list").
		Select("ip_list.*,man_server.id as msId,man_server.name as man_server_name, man_server.addr as addr,man_server.socks_port as socks_port,man_server.http_port as http_port,man_server.ssu_port as ssu_port").
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
		Select("user_project.*, ip_list.ip as ip,man_server.id as msId, man_server.socks_port as socks_port, man_server.http_port as http_port, man_server.ssu_port as ssu_port").
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
func AddShort(Info ShortUserProject) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Info).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func GetInfo(id int64) UserProject {
	var ret UserProject
	if err := Conn.Table("user_project").Where("id = ?", id).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetCity() []CITY {
	var results []CITY
	if err := Conn.Table("man_server").
		Select("man_server.id as id, man_server.name as name, city.id as CityId, city.name as CityName").
		Joins("join city on man_server.city_id = city.id").
		Scan(&results).Error; err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return results
}

func GetCitys() []City {
	ret := make([]City, 0)
	if err := Conn.Table("city").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret

}
func GetPr() []ManServer {
	ret := make([]ManServer, 0)
	if err := Conn.Table("man_server").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret

}
func MergeData(manServers []ManServer, cities []City) []CITY {
	cityMap := make(map[int]City)
	for _, city := range cities {
		cityMap[city.Id] = city
	}

	var combinedResults []CITY
	for _, server := range manServers {
		if city, ok := cityMap[server.CityId]; ok {
			combined := CITY{
				MsId:   server.Id,
				MsName: server.Name,
				//CityId:   city.Id,
				CityName: city.Name,
			}
			combinedResults = append(combinedResults, combined)
		}
	}

	return combinedResults
}
func Getsss(ids string) []IpList {
	var ret []IpList

	idSlice := strings.Split(ids, ",")
	var idIntSlice []int

	for _, idStr := range idSlice {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("无法将字符串 %s 转换为整数: %s\n", idStr, err.Error())
			return ret
		}
		idIntSlice = append(idIntSlice, id)
	}
	if err := Conn.Table("ip_list").Where("ms_i_d_id IN (?)", idIntSlice).Find(&ret).Error; err != nil {
		fmt.Printf("错误：%s\n", err.Error())
	}

	return ret
}
