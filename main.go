package main

import (
	"log"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"net/http"
	"time"
)
var Db *gorm.DB


func main() {
	if os.Getenv("DEP" ) == "DEP"{
		time.Sleep(time.Second * 5)
	}
 	e := echo.New()

	DbUser := os.Getenv("MYSQL_USER")
	DbPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	DbURL := os.Getenv("MYSQL_URL")
	DbName := os.Getenv("MYSQL_DATABASE")

	log.Println("MySQL Connectiong...")
	db, err := gorm.Open("mysql", DbUser+":"+DbPass+"@"+"tcp("+DbURL+")/"+DbName+"?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Error: Mysql conntection found")
		panic(err)
	}

	if !db.HasTable(&AccessLog{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").CreateTable(&AccessLog{})
	}
	Db = db

	e.GET("/",Index)

	e.Start(":" + os.Getenv("PORT"))
}

func Index(e echo.Context) error{
	var accessLog AccessLog
	accessLog.Address = e.RealIP()
	accessLog.Time = time.Now().Unix()
	log.Println(accessLog)
	Db.Create(&accessLog)
	if accessLog.Id == 0{
		return e.JSON(http.StatusInternalServerError,nil)
	}
	return e.JSON(http.StatusOK,nil)
}

type AccessLog struct{
	Id    int    `json:"Id",gorm:"primary_key",gorm:"AUTO_INCREMENT"`
	Address string `json:"Address" sql:"type:text"`
	Time     int64  `json:"Time"`
}
