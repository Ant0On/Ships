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
	_, err := cl.InitGame("You can run, but you can't hide", "Ant0n", true)
	if err != nil {
		log.Println(err)
	}
	cords, boardErr := cl.Board()
	if boardErr != nil {
		log.Println(boardErr)
	}

	desc, descErr := cl.Descriptions()
	if descErr != nil {
		log.Println(descErr)
	}
	fmt.Println(desc.Opponent)

	ui := gui.NewGUI(true)
	myBoard := gui.NewBoard(1, 3, nil)
	wpBoard := gui.NewBoard(50, 3, nil)
	ui.Draw(myBoard)
	ui.Draw(wpBoard)
	ui.Draw(gui.NewText(0, 0, desc.Nick, nil))
	ui.Draw(gui.NewText(50, 0, desc.Opponent, nil))
	ui.Draw(gui.NewText(0, 25, desc.Desc, nil))
	ui.Draw(gui.NewText(50, 25, desc.OppDesc, nil))
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
	//	ui.Start(nil)

	i := 0
	status, statsErr := cl.Status()
	for {
		status, statsErr = cl.Status()
		time.Sleep(time.Second * 1)
		fmt.Println(status)
		if statsErr != nil {
			return
		}
		if status.ShouldFire == true {
			for {
				var coord string
				fmt.Println("It's your turn!")
				fmt.Print("Enter coord: ")
				fmt.Scanf("%s", &coord)
				fire, fireErr := cl.Fire(coord)
				fmt.Println(fire)
				if fireErr != nil {
					log.Println(fireErr)
				}
				myBoard.SetStates(states)
				if fire == "miss" {
					break
				}
				status, _ = cl.Status()
			}
			status, _ = cl.Status()
			myBoard.SetStates(states)
		}
		if i == 5 {
			break
		}
		i++
	}

}
