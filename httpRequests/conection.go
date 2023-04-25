package httpRequests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	urlInit  = "https://go-pjatk-server.fly.dev/api/game"
	urlBoard = "https://go-pjatk-server.fly.dev/api/game/board"
)

func InitGame() error {

	type Data struct {
		Coords     []string `json:"coords"`
		Desc       string   `json:"desc"`
		Nick       string   `json:"nick"`
		TargetNick string   `json:"target_nick"`
		Wpbot      bool     `json:"wpbot"`
	}

	data := &Data{Wpbot: true}
	d, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}

	resp, err := http.Post(urlInit, "application/json", bytes.NewReader(d))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp.Header.Get("X-Auth-Token"))
	return err
}

func Board() ([]string, error) {
	type Data struct {
		Board []string `json:"board"`
	}
	resp, err := http.Get(urlBoard)
	if err != nil {
		log.Println(err)
	}
	var data Data
	err2 := json.NewDecoder(resp.Body).Decode(&data)
	if err2 != nil {
		log.Print(err2)
	}
	return data.Board, err2
}
