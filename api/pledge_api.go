package main

import "pledge-backend-study/db"

func main() {
	//先初始化mysql
	db.InitMysql()
	// 启动服务
	Init()
}
