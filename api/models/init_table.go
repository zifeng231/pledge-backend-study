package models

import "pledge-backend-study/db"

func InitTable() {
	db.Mysql.AutoMigrate(&MultiSign{})
	db.Mysql.AutoMigrate(&TokenInfo{})
	db.Mysql.AutoMigrate(&TokenList{})
	db.Mysql.AutoMigrate(&PoolData{})
	db.Mysql.AutoMigrate(&PoolBases{})
}
