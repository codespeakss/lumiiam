package main

import (
	"github.com/gin-gonic/gin"
	_ "go.uber.org/automaxprocs"
	"log"
	"lumiiam/config"
	"lumiiam/internal/handler"
	"lumiiam/internal/model"
	"lumiiam/pkg/cache"
	"lumiiam/pkg/db"
	"runtime"
)

func init() {
	PrintVersion()

	//var e error

	log.Println("Current GOMAXPROCS: ", runtime.GOMAXPROCS(0))
}

const Port = ":8011"

func main() {
	//	initMySQL & redis
	dbConn, err := db.Connect(config.MySQLDSN)
	if err != nil {
		log.Fatalf("db.Connect(config.MySQLDSN): %s", err)
	}
	log.Println("dbConn Connect OK")
	if err := dbConn.Migrator().DropTable(
		&model.User{},
	); err != nil {
		log.Fatalf("failed to DropTable: %v", err)
	}
	if err := dbConn.AutoMigrate(
		&model.User{},
	); err != nil {
		log.Fatalf("failed to auto migrate database: %v", err)
	}
	redis := cache.NewRedisTokenStore(config.RedisDsn)

	// http server
	r := gin.Default()
	h := handler.NewHandler(dbConn, redis)
	h.RegisterRoutes(r)
	r.Run(Port)

	log.Println("waiting select {}")
	select {}
}
