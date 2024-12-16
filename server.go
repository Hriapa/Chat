package main

import (
	"log"
	"webServer/apiserver"
	"webServer/store"

	"github.com/gorilla/sessions"
)

func main() {

	//var err error

	s := store.NewDataBase()
	s.Config.ReadConfig("store_config.json")
	if err := s.Open(); err != nil {
		log.Println("error connect to database", err.Error())
	} else {
		log.Println("database is connected")
	}

	ses := sessions.NewCookieStore(apiserver.ReadSequreParam("sequre.json"))

	server := apiserver.New(s, ses)

	server.Config.ReadFile("server_config.json")

	log.Fatal(server.Start())
}
