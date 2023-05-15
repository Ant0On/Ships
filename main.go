package main

import (
	"Ships/client"
	"Ships/logic"
	"log"
	"time"
)

const (
	baseURL = "https://go-pjatk-server.fly.dev/api"
	timeOut = time.Minute
)

func main() {

	appClient := client.NewClient(baseURL, timeOut)
	app := logic.NewApp(appClient)
	appErr := app.Run()
	if appErr != nil {
		log.Println(appErr)
		return
	}

}
