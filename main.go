package main

import (
	"encoding/json"
	"fmt"
	"record"
)

func main() {
	var usr1 record.User
	var usr2 record.User
	var usr3 record.User
	var usr4 record.User

	var usr5 record.User
	var usr6 record.User
	var usr7 record.User
	var usr8 record.User

	usr1.UserId = 1
	usr1.UserName = "player1"
	usr1.UserImg = "./img1"
	//record.CreateUserInfo(usr1)

	usr2.UserId = 2
	usr2.UserName = "player2"
	usr2.UserImg = "./img2"
	//record.CreateUserInfo(usr2)

	usr3.UserId = 3
	usr3.UserName = "player3"
	usr3.UserImg = "./img3"
	//record.CreateUserInfo(usr3)

	usr4.UserId = 4
	usr4.UserName = "player4"
	usr4.UserImg = "./img4"
	//record.CreateUserInfo(usr4)

	usr5.UserId = 5
	usr5.UserName = "player5"
	usr5.UserImg = "./img5"
	//record.CreateUserInfo(usr5)

	usr6.UserId = 6
	usr6.UserName = "player6"
	usr6.UserImg = "./img6"
	//record.CreateUserInfo(usr6)

	usr7.UserId = 7
	usr7.UserName = "player7"
	usr7.UserImg = "./img7"
	//record.CreateUserInfo(usr7)

	usr8.UserId = 8
	usr8.UserName = "player8"
	usr8.UserImg = "./img8"
	//record.CreateUserInfo(usr8)

	var r record.Record
	var d record.Detail
	// record.GameStart(usr1, usr2, usr3, usr4)
	// record.GameStart(usr5, usr6, usr7, usr8)
	// record.GameStart(usr1, usr3, usr5, usr7)
	data := record.GetRecordInfo(usr1)
	err := json.Unmarshal([]byte(data), &r)
	if err != nil {
		fmt.Printf("format err:%s\n", err.Error())
		return
	}
	fmt.Println(r)
	fmt.Println("--------------------------------------")

	data = record.GetGameDetail(usr1, 1)
	err = json.Unmarshal([]byte(data), &d)
	if err != nil {
		fmt.Printf("format err:%s\n", err.Error())
		return
	}
	fmt.Println(d)

}
