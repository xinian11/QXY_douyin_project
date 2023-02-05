package console

import (
	"database/sql"

	//go get -u gitub.com/go-sql-driver/mysql
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

//init database douyin_db
func InitDB() {
	dsn := "root:1190302927@tcp(127.0.0.1:3306)/douyin_db?charset=utf8mb4&parseTime=True"
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}
