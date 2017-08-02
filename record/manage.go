package record

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

//数据库相关信息
var (
	DBserverName = "root"
	Password     = "mysql"
	HostIP       = "127.0.0.1:3306"
	DBName       = "mahjong"
)

var M *DbManager
var Once sync.Once

//数据库连接字符串
var DbServer string = DBserverName + ":" + Password + "@tcp(" + HostIP + ")/" + DBName + "?charset=utf8"

//数据库管理
type DbManager struct {
	Db *sql.DB
}

//用户信息
type User struct {
	UserId   int    //用户id
	UserName string //用户名
	UserImg  string //用户头像
}

//玩家信息
type Player struct {
	PlayerName string //玩家名
	PlayerImag string //玩家头像
	AGameScore int    //一局游戏总积分
}

//用户单局战绩信息
type Score struct {
	UserId int //参与玩家用户Ｉｄ
	Score  int //单局积分
}

//战绩详情
type RecordInfo struct {
	GameId    int       //游戏id
	GameCode  int       //游戏编号
	BeginTime time.Time //游戏开始时间
	UserID    int       //当前玩家id
	WinUser   int       //获胜玩家id
	User1     Score     //参与玩家１
	User2     Score     //参与玩家2
	User3     Score     //参与玩家3
	User4     Score     //参与玩家4
}

//战绩统计
type StatisInfo struct {
	UserId     int     //用户id
	FinalScore int     //总积分
	TotalGames int     //总局数
	WinOfRate  float64 //胜率
}

//一大局游戏总计信息
type AGame struct {
	GameTime string   //游戏开始时间
	Total    int      //小局数
	Players  []Player //玩家战绩信息
}

//返回查询战绩信息
type Record struct {
	CurrentScore int     //当前积分
	TotalGames   int     //总局数
	WinOfRate    float64 //胜率
	GameDetail   []AGame
}

//大局游戏详情信息
type Detail struct {
	GameTime   string  //游戏开始时间
	Total      int     //小局数
	GameDetail []SGame //玩家战绩信息
}

//每小局的游戏战绩信息
type SGame struct {
	GameCode int
	Players  []Player
}
