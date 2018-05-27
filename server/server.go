package server

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

type Server struct {
	DB *gorm.DB
}

func (i *Server) InitDB() {
	var err error
	i.DB, err = gorm.Open("mysql", "root@/aruba?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	i.DB.LogMode(true)
}
