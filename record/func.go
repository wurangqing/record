/* 模拟获取战绩相关信息
------１．获取当前最近１００场游戏信息
------２．查看每场游戏的得分详情
*/

package record

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

//用于计算每大局游戏的自增量
var gameId int

//单例模式访问数据库
func GetInstance() *DbManager {
	Once.Do(func() {
		M = &DbManager{}
	})
	return M
}

//连接数据库
func (M *DbManager) DbConnect() *sql.DB {
	db1 := M.Db
	db1, err := sql.Open("mysql", DbServer)
	checkErr(err)
	return db1
}

//插入用户到用户表中
func CreateUserInfo(user User) {
	GetInstance()
	db := M.DbConnect()
	stmt, err := db.Prepare("insert user_info set user_id=?,user_name=?,user_img=?")
	checkErr(err)
	res, err := stmt.Exec(user.UserId, user.UserName, user.UserImg)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect, "!!!!!!insert a user user_info success!!!!!!")

}

//一盘游戏开始
func GameStart(user1 User, user2 User, user3 User, user4 User) {
	//判断进入游戏的四个用户在数据库中的游戏记录是否超过１００场，如果超过删除最早的，保持数据库中每个玩家存在１００场最近游戏记录
	RefreshRecord(user1)
	RefreshRecord(user2)
	RefreshRecord(user3)
	RefreshRecord(user4)
	//一大局游戏标识（唯一，自增）
	gameId++
	//获取数据库实例句柄
	GetInstance()
	//连接数据库
	db := M.DbConnect()

	var ri RecordInfo
	var si1, si2, si3, si4 StatisInfo

	//判断游戏中的四个玩家在战绩统计表中是否有记录，如果没有，就插入没有的玩家ＩＤ至战绩统计表中，如果存在，就不做任何操作
	if !RcordIsExist(user1) {
		InsertStatisInfoUserId(user1)
	}
	if !RcordIsExist(user2) {
		InsertStatisInfoUserId(user2)
	}
	if !RcordIsExist(user3) {
		InsertStatisInfoUserId(user3)
	}
	if !RcordIsExist(user4) {
		InsertStatisInfoUserId(user4)
	}

	//初始化信息
	si1.UserId = user1.UserId
	si2.UserId = user2.UserId
	si3.UserId = user3.UserId
	si4.UserId = user4.UserId

	ri.GameId = gameId
	ri.BeginTime = time.Now()

	ri.User1.UserId = user1.UserId
	ri.User2.UserId = user2.UserId
	ri.User3.UserId = user3.UserId
	ri.User4.UserId = user4.UserId

	//游戏循环八小局
	for i := 1; i < 9; i++ {
		ri.GameCode = i
		//随机种子
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		//随机分配分数给四个玩家（模拟游戏得分）
		ri.User1.Score = r.Intn(100) * i
		ri.User2.Score = r.Intn(100) * i
		ri.User3.Score = r.Intn(100) * i
		ri.User4.Score = r.Intn(100) * i
		//计算本局最高得分并获取其玩家ＩＤ赋值给获胜玩家ＩＤ
		WinScore := WinUser(ri.User1.Score, ri.User2.Score, ri.User3.Score, ri.User4.Score)
		switch WinScore {
		case ri.User1.Score:
			ri.WinUser = user1.UserId
		case ri.User2.Score:
			ri.WinUser = user2.UserId
		case ri.User3.Score:
			ri.WinUser = user3.UserId
		case ri.User4.Score:
			ri.WinUser = user4.UserId
		default:
			fmt.Println("------error------")
		}
		//插入四个玩家一局游戏战绩信息到战绩详情表中
		stmt, err := db.Prepare("INSERT record_info SET game_id=?,game_code=?,begin_time=?,user_id=?,win_user=?,join_user1=?,score1=?,join_user2=?,score2=?,join_user3=?,score3=?,join_user4=?,score4=?")
		checkErr(err)
		res, err := stmt.Exec(ri.GameId, ri.GameCode, ri.BeginTime, user1.UserId, ri.WinUser, ri.User1.UserId, ri.User1.Score, ri.User2.UserId, ri.User2.Score, ri.User3.UserId, ri.User3.Score, ri.User4.UserId, ri.User4.Score)
		res, err = stmt.Exec(ri.GameId, ri.GameCode, ri.BeginTime, user2.UserId, ri.WinUser, ri.User1.UserId, ri.User1.Score, ri.User2.UserId, ri.User2.Score, ri.User3.UserId, ri.User3.Score, ri.User4.UserId, ri.User4.Score)
		res, err = stmt.Exec(ri.GameId, ri.GameCode, ri.BeginTime, user3.UserId, ri.WinUser, ri.User1.UserId, ri.User1.Score, ri.User2.UserId, ri.User2.Score, ri.User3.UserId, ri.User3.Score, ri.User4.UserId, ri.User4.Score)
		res, err = stmt.Exec(ri.GameId, ri.GameCode, ri.BeginTime, user4.UserId, ri.WinUser, ri.User1.UserId, ri.User1.Score, ri.User2.UserId, ri.User2.Score, ri.User3.UserId, ri.User3.Score, ri.User4.UserId, ri.User4.Score)
		checkErr(err)
		affect, err := res.RowsAffected()
		checkErr(err)
		fmt.Println(affect, "!!!!!!insert a game record_info success!!!!!!")
		//统计四个玩家的总分和总局数
		si1.FinalScore = ri.User1.Score
		si2.FinalScore = ri.User2.Score
		si3.FinalScore = ri.User3.Score
		si4.FinalScore = ri.User4.Score

		si1.TotalGames = 1
		si2.TotalGames = 1
		si3.TotalGames = 1
		si4.TotalGames = 1

		//更新战绩统计表的信息
		stmt, err = db.Prepare("update statis_info SET final_score=final_score+?,total_games=total_games+? WHERE user_id=?")
		checkErr(err)
		res, err = stmt.Exec(si1.FinalScore, si1.TotalGames, si1.UserId)
		checkErr(err)

		stmt, err = db.Prepare("UPDATE statis_info SET final_score=final_score+?,total_games=total_games+?  WHERE user_id=?")
		checkErr(err)
		res, err = stmt.Exec(si2.FinalScore, si2.TotalGames, si2.UserId)
		checkErr(err)

		stmt, err = db.Prepare(`UPDATE statis_info SET final_score=final_score+?,total_games=total_games+? WHERE user_id=?`)
		checkErr(err)
		res, err = stmt.Exec(si3.FinalScore, si3.TotalGames, si3.UserId)
		checkErr(err)

		stmt, err = db.Prepare(`UPDATE statis_info SET final_score=final_score+?,total_games=total_games+? WHERE user_id=?`)
		checkErr(err)
		res, err = stmt.Exec(si4.FinalScore, si4.TotalGames, si4.UserId)
		checkErr(err)
		num, err := res.RowsAffected()
		checkErr(err)
		fmt.Println(num)
		fmt.Println(num, "!!!!!!update statis_info success!!!!!!")

	}
	//更新玩家的胜率
	si1.WinOfRate = float64(CalcWinOfRate(user1)) / float64(FindUserTotalGames(user1))
	si2.WinOfRate = float64(CalcWinOfRate(user2)) / float64(FindUserTotalGames(user2))
	si3.WinOfRate = float64(CalcWinOfRate(user3)) / float64(FindUserTotalGames(user3))
	si4.WinOfRate = float64(CalcWinOfRate(user4)) / float64(FindUserTotalGames(user4))

	stmt, err := db.Prepare("update statis_info SET win_of_rate=? WHERE user_id=?")
	checkErr(err)
	res, err := stmt.Exec(si1.WinOfRate, si1.UserId)
	checkErr(err)

	stmt, err = db.Prepare("UPDATE statis_info SET win_of_rate=?  WHERE user_id=?")
	checkErr(err)
	res, err = stmt.Exec(si2.WinOfRate, si2.UserId)
	checkErr(err)

	stmt, err = db.Prepare(`UPDATE statis_info SET win_of_rate=?  WHERE user_id=?`)
	checkErr(err)
	res, err = stmt.Exec(si3.WinOfRate, si3.UserId)
	checkErr(err)

	stmt, err = db.Prepare(`UPDATE statis_info SET win_of_rate=?  WHERE user_id=?`)
	checkErr(err)
	res, err = stmt.Exec(si4.WinOfRate, si4.UserId)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)
	fmt.Println(num, "!!!!!!update statis_info success!!!!!!")

	//关闭数据库连接
	defer db.Close()

}

//求用户获胜的的总局数
func CalcWinOfRate(user User) int {
	GetInstance()
	var count int
	db := M.DbConnect()
	rows, err := db.Query("select count(*) from record_info where win_user=?", user.UserId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&count)
	}
	return count / 4
}

//查询用户的总局数
func FindUserTotalGames(user User) int {
	GetInstance()
	var count int
	db := M.DbConnect()
	rows, err := db.Query("select total_games from statis_info where user_id=?", user.UserId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&count)
	}
	return count
}

//查看该用户在统计战绩表中是否有记录
func RcordIsExist(user User) bool {
	GetInstance()
	db := M.DbConnect()
	var count int
	var x bool
	rows, err := db.Query("select count(*) from statis_info where user_id=?", user.UserId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&count)
	}
	if count < 1 {
		x = false
	} else {
		x = true
	}
	return x
}

//插入战绩统计表中不存在的用户的ＩＤ
func InsertStatisInfoUserId(user User) {
	GetInstance()
	db := M.DbConnect()
	stmt, err := db.Prepare("INSERT statis_info SET user_id=?")
	checkErr(err)
	res, err := stmt.Exec(user.UserId)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)

}

//求出赢得一小局游戏的赢家
func WinUser(score1, score2, score3, score4 int) int {
	var max, max1, max2 int
	if score1 < score2 {
		max1 = score2
	} else {
		max1 = score1
	}
	if score3 < score4 {
		max2 = score4
	} else {
		max2 = score3
	}
	if max1 < max2 {
		max = max2
	} else {
		max = max1
	}
	return max
}

//更新数据记录，保持每个玩家的信息记录只有１００条
func RefreshRecord(user User) {
	GetInstance()
	db := M.DbConnect()

	var count int
	rows, err := db.Query("SELECT count(*) FROM record_info where user_id=?", user.UserId)
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&count)
		checkErr(err)
	}

	if count > 100 {
		//删除多余的数据库记录
		stmt, err := db.Prepare("delete from record_info where user_id=? and begin_time like (select begin_time from(select min(begin_time)as begin_time from record_info where user_id=?)as temp)")
		checkErr(err)
		res, err := stmt.Exec(user.UserId, user.UserId)
		checkErr(err)
		num, err := res.RowsAffected()
		checkErr(err)
		fmt.Println(num)
	}
	defer db.Close()
}

//查询战绩信息（最近100场）
func GetRecordInfo(user User) string {
	//获取个人战绩统计情况
	var record Record

	var s StatisInfo
	s = FindUserStatisInfo(user.UserId)
	record.CurrentScore = s.FinalScore
	record.TotalGames = s.TotalGames
	record.WinOfRate = s.WinOfRate
	count := FindUserGameId(user.UserId)
	//-----------------------
	//fmt.Println(record, count)
	var gameId int
	for _, i := range count {

		gameId = i
		//fmt.Println(gameId)
		//fmt.Println("```````````````", user.UserId, gameId)
		g.GameTime = FindGameTime(user.UserId, gameId)
		g.Total = 8
		var m map[string]int = make(map[string]int, 4)
		m = FindPlayerId(user.UserId, gameId)

		var p1, p2, p3, p4 Player
		p1.PlayerName = FindUserInfo(m["join_user1"]).UserName
		p1.PlayerImag = FindUserInfo(m["join_user1"]).UserImg
		p1.AGameScore = CalcGameScore(user.UserId, gameId, m["join_user1"])
		g.Players = append(g.Players, p1)
		fmt.Println("////", p1)

		p2.PlayerName = FindUserInfo(m["join_user2"]).UserName
		p2.PlayerImag = FindUserInfo(m["join_user2"]).UserImg
		p2.AGameScore = CalcGameScore(user.UserId, gameId, m["join_user2"])
		g.Players = append(g.Players, p2)

		p3.PlayerName = FindUserInfo(m["join_user3"]).UserName
		p3.PlayerImag = FindUserInfo(m["join_user3"]).UserImg
		p3.AGameScore = CalcGameScore(user.UserId, gameId, m["join_user3"])
		g.Players = append(g.Players, p3)

		p4.PlayerName = FindUserInfo(m["join_user4"]).UserName
		p4.PlayerImag = FindUserInfo(m["join_user4"]).UserImg
		p4.AGameScore = CalcGameScore(user.UserId, gameId, m["join_user4"])
		g.Players = append(g.Players, p4)

		record.GameDetail = append(record.GameDetail, g)
	}
	bytes, _ := json.Marshal(record)
	send := string(bytes)
	return send
}

//查询每局游戏详情（８小局游戏战绩情况）
func GetGameDetail(user User, gameId int) string {
	var d Detail
	d.GameTime = FindGameTime(user.UserId, gameId)
	d.Total = 8

	var s1, s2, s3, s4 int
	GetInstance()
	db := M.DbConnect()
	for i := 1; i < 9; i++ {
		var sg SGame
		gameCode := i
		var m map[string]int = make(map[string]int, 4)

		m = FindPlayerId(user.UserId, gameId)

		rows, err := db.Query("select game_code,score1,score2,score3,score4 from record_info where user_id=? and game_id=? and game_code=?", user.UserId, gameId, gameCode)
		checkErr(err)
		for rows.Next() {
			rows.Scan(&sg.GameCode, &s1, &s2, &s3, &s4)
		}

		var p1, p2, p3, p4 Player
		p1.PlayerName = FindUserInfo(m["join_user1"]).UserName
		p1.PlayerImag = FindUserInfo(m["join_user1"]).UserImg
		//fmt.Println("////------", s1)
		p1.AGameScore = s1
		sg.Players = append(sg.Players, p1)

		p2.PlayerName = FindUserInfo(m["join_user2"]).UserName
		p2.PlayerImag = FindUserInfo(m["join_user2"]).UserImg
		p2.AGameScore = s2
		sg.Players = append(sg.Players, p2)
		//fmt.Println("////", p2)

		p3.PlayerName = FindUserInfo(m["join_user3"]).UserName
		p3.PlayerImag = FindUserInfo(m["join_user3"]).UserImg
		p3.AGameScore = s3
		sg.Players = append(sg.Players, p3)

		p4.PlayerName = FindUserInfo(m["join_user4"]).UserName
		p4.PlayerImag = FindUserInfo(m["join_user4"]).UserImg
		p4.AGameScore = s4
		sg.Players = append(sg.Players, p4)

		//fmt.Println("-----", sg)
		// var x int
		// x = i - 1
		d.GameDetail = append(d.GameDetail, sg)

	}

	bytes, _ := json.Marshal(d)
	send := string(bytes)
	return send
}

//错误处理
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//查询用户信息
func FindUserInfo(userId int) User {
	var result User
	GetInstance()
	//var count int
	db := M.DbConnect()
	rows, err := db.Query("select * from user_info where user_id=?", userId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&result.UserId, &result.UserName, &result.UserImg)
	}
	return result

}

//查询用户的统计战绩信息
func FindUserStatisInfo(userId int) StatisInfo {
	var result StatisInfo
	GetInstance()
	//var count int
	db := M.DbConnect()
	rows, err := db.Query("select * from statis_info where user_id=?", userId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&result.UserId, &result.FinalScore, &result.TotalGames, &result.WinOfRate)
	}
	return result
}

//查询某一用户所参加最近１００场的游戏大局编号
func FindUserGameId(userId int) []int {
	var result []int
	var t int
	result = make([]int, 0)
	GetInstance()
	db := M.DbConnect()
	rows, err := db.Query("select distinct game_id from record_info where user_id=?", userId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&t)
		result = append(result, t)
	}
	return result
}

//根据用户id和游戏id查询该局游戏所有玩家id
func FindPlayerId(userId, gameId int) map[string]int {
	//fmt.Println("--------------", gameId, userId)
	var result map[string]int
	result = make(map[string]int)
	var t1, t2, t3, t4 int
	GetInstance()
	db := M.DbConnect()
	rows, err := db.Query("select distinct join_user1,join_user2,join_user3,join_user4 from record_info where game_id=? and user_id=?", gameId, userId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&t1, &t2, &t3, &t4)
	}
	if t1 != 0 && t2 != 0 && t3 != 0 && t4 != 0 {
		result["join_user1"] = t1
		result["join_user2"] = t2
		result["join_user3"] = t3
		result["join_user4"] = t4

		//fmt.Println(result)
		return result

	} else {
		result["error"] = 0
		return result
	}

}

//根据用户id玩家id游戏id计算每大局得分
func CalcGameScore(userId, gameId, playerId int) int {
	pid := FindPlayerId(userId, gameId)
	// var pid map[string]int = make(map[string]int)
	// if ok {
	// 	pid = p
	// }
	var s, flag string
	//fmt.Println("///////////", pid)
	for k, v := range pid {
		if playerId == v {
			//fmt.Println(k, v)
			s = k
			break
		}
	}
	//fmt.Println("********", s)
	switch s {
	case "join_user1":
		flag = "score1"
	case "join_user2":
		flag = "score2"
	case "join_user3":
		flag = "score3"
	case "join_user4":
		flag = "score4"
	default:
		fmt.Println("sorry,no record!")

	}

	str := "select sum(" + flag + ")" + " from record_info where game_id=? and user_id=? and " + s + "=?"
	//fmt.Println(flag, str)
	var result int
	GetInstance()
	db := M.DbConnect()
	rows, err := db.Query(str, gameId, userId, playerId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&result)
	}
	return result
}

//查询某个玩家某局游戏开始时间
func FindGameTime(userId, gameId int) string {
	var t string
	GetInstance()
	db := M.DbConnect()
	rows, err := db.Query("select distinct begin_time from record_info where game_id=? and user_id=? ", gameId, userId)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&t)
	}
	//fmt.Println(t)
	return t
}
