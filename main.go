package main

import (
	"bible_reading_backend_nkv/database"
	"bible_reading_backend_nkv/server"
	"log"
)

func main (){
	db, err:= database.NewDatabaseClient()
	if err!=nil {
		log.Fatalf("failed to initialize Database Client: %s", err)
	}
	serv:=server.NewEchoServer(db)
	if err:=serv.Start(); err!=nil{
		log.Fatal(err.Error())

	}
	log.Fatal("ash")
}