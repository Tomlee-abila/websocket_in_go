package main

import (
	"log"
	"net/http"
	"websocket/utills"
)

func main(){
	setupAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupAPI(){
	manager := utills.NewManager();

	http.Handle("/", http.FileServer(http.Dir("./frontend")));
	http.HandleFunc("/ws", manager.ServeWs)
}

