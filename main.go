package main

import (
	"Ships/client"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"time"
)

const (
	baseURL = "https://go-pjatk-server.fly.dev/api"
	timeOut = time.Minute
)

func main() {

	cl := client.NewClient(baseURL, timeOut)
	err := cl.InitGame()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(cl.Token)
	cords, _ := cl.Board()

	_, stats, statsErr := cl.Status()
	if statsErr != nil {
		log.Println(statsErr)
	}
	fmt.Println(stats)

	ui := gui.NewGUI(true)
	myBoard := gui.NewBoard(1, 3, nil)
	wpBoard := gui.NewBoard(50, 3, nil)
	ui.Draw(myBoard)
	ui.Draw(wpBoard)
	ui.Draw(gui.NewText(0, 0, "Antoni", nil))
	ui.Draw(gui.NewText(50, 0, stats.Opponent, nil))
	states := [10][10]gui.State{}
	for i := range states {
		states[i] = [10]gui.State{}
	}
	stateSlice := make([][]gui.State, len(states))
	for i := range stateSlice {
		stateSlice[i] = states[i][:]
	}

	client.PlaceShips(cords, stateSlice)
	myBoard.SetStates(states)
	ui.Start(nil)

}
