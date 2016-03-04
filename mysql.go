package gwork

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func MysqlQuery(sqlString string) {
	db, err := sql.Open("mysql", "root:mike0125@/dixit?charset=utf8")
	CheckErr(err)
	_, err = db.Query(sqlString)
	CheckErr(err)
	db.Close()
}
