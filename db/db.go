package db

import (
	"database/sql"
	"log"

	"user-bank-manage/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func Init() {
	conf := config.GetConfig()

	// connectionString := "root:Padaringan_k383k@tcp(127.0.0.1:3306)/localpedia"
	connectionString := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		panic("connectionString error..")
	}

	err = db.Ping()
	if err != nil {
		// panic(err)
		log.Println(conf.DB_USERNAME)
		log.Println(conf.DB_PASSWORD)
		log.Println(conf.DB_HOST)
		log.Println(conf.DB_PORT)
		log.Println(conf.DB_NAME)
	}
}

func CreateCon() *sql.DB {
	return db
}
