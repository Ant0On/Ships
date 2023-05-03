package main

import (
	"Ships/client"
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
	initialData, err := cl.InitGame("You can run, but you can't hide", "Ant0n", true)
	if err != nil {
		log.Println(err)
	}
	cords, boardErr := cl.Board()
	if boardErr != nil {
		log.Println(boardErr)
	}
	time.Sleep(time.Second * 1)
	stats, statsErr := cl.Status()
	if statsErr != nil {
		log.Println(statsErr)
	}
	ui := gui.NewGUI(true)
	myBoard := gui.NewBoard(1, 3, nil)
	wpBoard := gui.NewBoard(50, 3, nil)
	ui.Draw(myBoard)
	ui.Draw(wpBoard)
	ui.Draw(gui.NewText(0, 0, initialData.Nick, nil))
	ui.Draw(gui.NewText(50, 0, stats.Opponent, nil))
	ui.Draw(gui.NewText(0, 25, initialData.Desc, nil))
	ui.Draw(gui.NewText(50, 25, stats.OppDesc, nil))
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
