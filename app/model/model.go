package model

import "time"

type IpList struct {
	Id      int    `gorm:"column:id;type:int(11);NOT NULL" json:"id"`
	Ip      string `gorm:"column:ip;type:varchar(150);NOT NULL" json:"ip"`
	MsIDId  int    `gorm:"column:ms_i_d_id;type:int(11);NOT NULL" json:"ms_i_d_id"`
	CityId  int    `gorm:"column:city_id;type:int(11);NOT NULL" json:"city_id"`
	Network string `gorm:"column:network;type:varchar(150);NOT NULL" json:"network"`
}

func (m *IpList) TableName() string {
	return "ip_list"
}

type ManServer struct {
	Id         int     `gorm:"column:id;type:int(11);NOT NULL" json:"id"`
	Name       string  `gorm:"column:name;type:varchar(60);NOT NULL" json:"name"`
	Addr       string  `gorm:"column:addr;type:varchar(150);NOT NULL" json:"addr"`
	SocksPort  int     `gorm:"column:socks_port;type:int(11);default:0;NOT NULL" json:"socks_port"`
	HttpPort   int     `gorm:"column:http_port;type:int(11);default:0;NOT NULL" json:"http_port"`
	SsuPort    int     `gorm:"column:ssu_port;type:int(11);default:0;NOT NULL" json:"ssu_port"`
	Money      float64 `gorm:"column:money;type:decimal(10,2);default:1.00;NOT NULL" json:"money"`
	Key        string  `gorm:"column:key;type:varchar(150);NOT NULL" json:"key"`
	Info       string  `gorm:"column:info;type:varchar(150);NOT NULL" json:"info"`
	Enable     int     `gorm:"column:enable;type:tinyint(4);default:1;NOT NULL" json:"enable"`
	CityId     int     `gorm:"column:city_id;type:int(11);NOT NULL" json:"city_id"`
	SocksState int     `gorm:"column:socks_state;type:tinyint(4);default:1;NOT NULL" json:"socks_state"`
	HttpState  int     `gorm:"column:http_state;type:tinyint(4);default:1;NOT NULL" json:"http_state"`
	SsuState   int     `gorm:"column:ssu_state;type:tinyint(4);default:1;NOT NULL" json:"ssu_state"`
	Ping       int     `gorm:"column:ping;type:tinyint(4);default:1;NOT NULL" json:"ping"`
}

func (m *ManServer) TableName() string {
	return "man_server"
}

type City struct {
	Id   int    `gorm:"column:id;type:int(11);NOT NULL" json:"id"`
	Name string `gorm:"column:name;type:varchar(30);NOT NULL" json:"name"`
}

func (m *City) TableName() string {
	return "city"
}

type UserProject struct {
	Id         int    `gorm:"column:id;type:int(11);NOT NULL" json:"id"`
	Type       int    `gorm:"column:type;type:int(11);default:1;NOT NULL" json:"type"`
	GameId     int    `gorm:"column:game_id;type:int(11);default:0;NOT NULL" json:"game_id"`
	UserId     int    `gorm:"column:user_id;type:int(11);default:0;NOT NULL" json:"user_id"`
	IpId       int    `gorm:"column:ip_id;type:int(11);default:0;NOT NULL" json:"ip_id"`
	Username   string `gorm:"column:username;type:varchar(30);NOT NULL" json:"username"`
	Password   string `gorm:"column:password;type:varchar(30);NOT NULL" json:"password"`
	MaxComm    int    `gorm:"column:max_comm;type:int(11);default:0;NOT NULL" json:"max_comm"`
	UpData     int    `gorm:"column:up_data;type:int(11);default:0;NOT NULL" json:"up_data"`
	DwData     int    `gorm:"column:dw_data;type:int(11);default:0;NOT NULL" json:"dw_data"`
	DeleteTime int    `gorm:"column:delete_time;type:int(11);default:0;NOT NULL" json:"delete_time"`
	CreatTime  int    `gorm:"column:creat_time;type:int(11);default:0;NOT NULL" json:"creat_time"`
	PutNode    int    `gorm:"column:put_node;type:int(11);default:1;NOT NULL" json:"put_node"`
	PutGame    int    `gorm:"column:put_game;type:int(11);default:1;NOT NULL" json:"put_game"`
	Gp         string `gorm:"column:gp;type:varchar(30);NOT NULL" json:"gp"`
}

func (m *UserProject) TableName() string {
	return "user_project"
}

type User struct {
	Id          int     `gorm:"column:id;type:int(11);NOT NULL" json:"id"`
	Username    string  `gorm:"column:username;type:varchar(30);NOT NULL" json:"username"`
	Mobile      string  `gorm:"column:mobile;type:varchar(11);NOT NULL" json:"mobile"`
	Password    string  `gorm:"column:password;type:varchar(255);NOT NULL" json:"password"`
	Status      int     `gorm:"column:status;type:tinyint(4);default:1;NOT NULL" json:"status"`
	Description string  `gorm:"column:description;type:longtext;NOT NULL" json:"description"`
	Wallet      float64 `gorm:"column:wallet;type:decimal(10,2);default:0.00;NOT NULL" json:"wallet"`
	CreateTime  int     `gorm:"column:create_time;type:int(11);default:0;NOT NULL" json:"create_time"`
	UserLevelId int     `gorm:"column:user_level_id;type:int(11);NOT NULL" json:"user_level_id"`
	ProFileId   int     `gorm:"column:pro_file_id;type:int(11);NOT NULL" json:"pro_file_id"`
	PopUser     int     `gorm:"column:pop_user;type:int(11);default:0;NOT NULL" json:"pop_user"`
}

func (m *User) TableName() string {
	return "user"
}

type ApiList struct {
	Id         int       `gorm:"column:id;type:int(11);NOT NULL" json:"id"`
	UserId     int       `gorm:"column:user_id;type:int(11);NOT NULL" json:"user_id"`
	IpIdList   string    `gorm:"column:ip_id_list;type:text;NOT NULL" json:"ip_id_list"`
	CreatTime  time.Time `gorm:"column:creat_time;type:datetime;NOT NULL" json:"creat_time"`
	DeleteTime time.Time `gorm:"column:delete_time;type:datetime;NOT NULL" json:"delete_time"`
}

func (m *ApiList) TableName() string {
	return "api_list"
}

type Game struct {
	Id      int     `gorm:"column:id;type:int(11);NOT NULL" json:"id"`
	Name    string  `gorm:"column:name;type:varchar(60);NOT NULL" json:"name"`
	MaxConn int     `gorm:"column:max_conn;type:int(11);default:0;NOT NULL" json:"max_conn"`
	UpData  int     `gorm:"column:up_data;type:int(11);default:0;NOT NULL" json:"up_data"`
	DwData  int     `gorm:"column:dw_data;type:int(11);default:0;NOT NULL" json:"dw_data"`
	Money   float64 `gorm:"column:money;type:decimal(10,2);default:1.00;NOT NULL" json:"money"`
	Info    string  `gorm:"column:info;type:varchar(150);NOT NULL" json:"info"`
	Info2   string  `gorm:"column:info2;type:varchar(500);NOT NULL" json:"info2"`
}

func (m *Game) TableName() string {
	return "game"
}

type ShortUserProject struct {
	Id         int    `gorm:"column:id;type:int(11);NOT NULL" json:"id"`
	Type       int    `gorm:"column:type;type:int(11);NOT NULL" json:"type"`
	GameId     int    `gorm:"column:game_id;type:int(11);NOT NULL" json:"game_id"`
	UserId     int    `gorm:"column:user_id;type:int(11);NOT NULL" json:"user_id"`
	BuyCount   int    `gorm:"column:buy_count;type:int(11);NOT NULL" json:"buy_count"`
	ModeId     int    `gorm:"column:mode_id;type:int(11);NOT NULL" json:"mode_id"`
	MsId       int    `gorm:"column:ms_id;type:int(11);NOT NULL" json:"ms_id"`
	DeleteTime string `gorm:"column:delete_time;type:date;NOT NULL" json:"delete_time"`
	CreatTime  string `gorm:"column:creat_time;type:date;NOT NULL" json:"creat_time"`
}

func (m *ShortUserProject) TableName() string {
	return "short_user_project"
}
