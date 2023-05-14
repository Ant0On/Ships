package main

import (
	"Ships/client"
	"Ships/logic"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"time"
)

const (
	baseURL = "https://go-pjatk-server.fly.dev/api"
	timeOut = time.Minute
	hit     = gui.Hit
	miss    = gui.Miss
)

func main() {

	//cl := client.NewClient(baseURL, timeOut)
	//_, err := cl.InitGame("You can run, but you can't hide", "Ant0n", true)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//ui := gui.NewGUI(true)
	//myBoard := gui.NewBoard(1, 3, nil)
	//wpBoard := gui.NewBoard(130, 3, nil)
	//go func() {
	//	for {
	//		status, statusErr := cl.Status()
	//		if statusErr != nil {
	//			log.Println(statusErr)
	//			return
	//		}
	//		cords, boardErr := cl.Board()
	//		if boardErr != nil {
	//			log.Println(boardErr)
	//			return
	//		}
	//
	//		desc, descErr := cl.Descriptions()
	//		if descErr != nil {
	//			log.Println(descErr)
	//			return
	//		}
	//
	//		ui.Draw(myBoard)
	//		ui.Draw(wpBoard)
	//		ui.Draw(gui.NewText(20, 1, desc.Nick, nil))
	//		ui.Draw(gui.NewText(140, 1, desc.Opponent, nil))
	//		ui.Draw(gui.NewText(20, 25, desc.Desc, nil))
	//		ui.Draw(gui.NewText(130, 25, desc.OppDesc, nil))
	//		states := [10][10]gui.State{}
	//		var opponentStates [10][10]gui.State
	//		for i := range states {
	//			states[i] = [10]gui.State{}
	//			opponentStates[i] = [10]gui.State{}
	//		}
	//		stateSlice := make([][]gui.State, len(states))
	//		for i := range stateSlice {
	//			stateSlice[i] = states[i][:]
	//		}
	//
	//		client.PlaceShips(cords, stateSlice)
	//		myBoard.SetStates(states)
	//		if status.GameStatus == "game_in_progress" {
	//
	//			txt := gui.NewText(70, 10, "Something", nil)
	//			ui.Draw(txt)
	//			status, statusErr = cl.Status()
	//			if statusErr != nil {
	//				log.Println(statusErr)
	//				return
	//			}
	//			if status.GameStatus == "ended" {
	//				break
	//			} else {
	//				for status.ShouldFire == true {
	//					if len(status.OppShots) != 0 {
	//						oppCoord := strings.Join(status.OppShots, "")
	//						x1, y1 := client.ConvertCords(oppCoord)
	//						if states[x1][y1] == gui.Hit {
	//							states[x1][y1] = gui.Hit
	//							myBoard.SetStates(states)
	//						} else {
	//							states[x1][y1] = gui.Miss
	//							myBoard.SetStates(states)
	//						}
	//					}
	//					txt.SetText("Your turn")
	//					time.Sleep(time.Second * 1)
	//					txt.SetText("Press on opponent's board to fire")
	//					coord := wpBoard.Listen(context.TODO())
	//					txt.SetText(fmt.Sprintf("Coordinate: %s", coord))
	//					fire, fireErr := cl.Fire(coord)
	//					if fireErr != nil {
	//						log.Println(fireErr)
	//						return
	//					}
	//					time.Sleep(time.Second * 1)
	//					if fire == "hit" || fire == "sunk" {
	//						txt.SetText(fire)
	//						x, y := client.ConvertCords(coord)
	//						opponentStates[x][y] = hit
	//						wpBoard.SetStates(opponentStates)
	//					} else {
	//						txt.SetText(fire)
	//						x, y := client.ConvertCords(coord)
	//						opponentStates[x][y] = miss
	//						wpBoard.SetStates(opponentStates)
	//						status.ShouldFire = false
	//					}
	//				}
	//				time.Sleep(time.Second * 1)
	//				status, _ = cl.Status()
	//			}
	//		} else {
	//			time.Sleep(time.Second * 1)
	//		}
	//	}
	//}()
	//ui.Start(nil)

	appClient := client.NewClient(baseURL, timeOut)
	app := logic.NewApp(appClient)
	appErr := app.Run()
	if appErr != nil {
		log.Println(appErr)
		return
	}

}
