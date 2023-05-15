package logic

import (
	"Ships/client"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"time"
)

type App struct {
	client *client.Client
}

var boardState BoardState

func NewApp(client *client.Client) *App {
	return &App{client: client}
}

func (app *App) Run() error {
	time.Sleep(time.Second * 1)
	listData, listDataErr := app.client.PlayersList()

	fmt.Println(listData)
	if listDataErr != nil {
		log.Println(listDataErr)
		return listDataErr
	}
	var targetNick string
	var nick string
	fmt.Print("Enter Your nick:")
	fmt.Scanln(&nick)
	fmt.Println("Enter enemy nick: ")
	fmt.Scanln(&targetNick)

	initErr := app.client.InitGame(nick, "You can run, but you can't hide", targetNick, false)
	if initErr != nil {
		log.Println(initErr)
		return initErr
	}

	startErr := app.waitToStart()
	if startErr != nil {
		log.Println(startErr)
		return startErr
	}
	description, descErr := app.client.Descriptions()
	if descErr != nil {
		log.Println(descErr)
		return descErr
	}
	myShips, shipsErr := app.client.Board()
	if shipsErr != nil {
		log.Println(shipsErr)
		return shipsErr
	}
	myBoard, enemyBoard, ui := CreateBoard(description)
	boardState.InitialStates(myBoard, enemyBoard, myShips)
	go func() error {
		appError := app.gameCourse(myBoard, enemyBoard, ui)
		if appError != nil {
			log.Println(appError)
			return appError
		}
		return nil
	}()
	ui.Start(nil)
	return nil
}

func (app *App) waitToStart() error {
	for {
		time.Sleep(time.Second * 1)
		gameStatus, statusErr := app.client.Status()
		if statusErr != nil {
			log.Println(statusErr)
			return statusErr
		}
		fmt.Println(gameStatus.GameStatus)
		if gameStatus.GameStatus == "game_in_progress" {
			break
		}
	}
	return nil
}

func (app *App) gameCourse(myBoard, enemyBoard *gui.Board, ui *gui.GUI) error {
	status, statusErr := app.client.Status()
	if statusErr != nil {
		log.Println(statusErr)
		return statusErr
	}
	hits := 0
	totalShoots := 0
	txt := gui.NewText(50, 5, "Let's start", nil)
	ui.Draw(txt)
	mark := gui.NewText(65, 15, "", nil)
	ui.Draw(mark)
	time.Sleep(time.Second * 2)
	for status.GameStatus == "game_in_progress" {
		for status.ShouldFire == true {
			if len(status.OppShots) != 0 {
				boardState.EnemyShoot(myBoard, status)
			}
			txt.SetBgColor(gui.White)
			txt.SetText("It's your turn!")
			time.Sleep(time.Second * 2)
			txt.SetText("Press on opponent's coordinates to shoot")
			var coords string
			for {
				coords = enemyBoard.Listen(context.TODO())
				x, y := ConvertCords(coords)
				if boardState.EnemyState[x][y] == "" {
					break
				}
			}
			boardState.AccuracyText = gui.NewText(5, 30, "AccuracyText: 0", nil)
			ui.Draw(boardState.AccuracyText)
			mark.SetBgColor(gui.White)
			mark.SetText(fmt.Sprintf("Coordinate: %s", coords))
			time.Sleep(time.Second * 1)

			fire, fireErr := app.client.Fire(coords)
			if fireErr != nil {
				log.Println(fireErr)
				return fireErr
			}
			if fire != "miss" {
				hits++
			}
			totalShoots++
			boardState.Accuracy = boardState.countAccuracy(hits, totalShoots)
			boardState.MarkMyShoot(mark, enemyBoard, fire, coords)
			status, _ = app.client.Status()
			time.Sleep(time.Second * 1)
		}
		time.Sleep(time.Second * 1)
		status, _ = app.client.Status()
	}
	if len(status.OppShots) != 0 {
		boardState.EnemyShoot(myBoard, status)
	}
	checkWinner(status, txt)
	return nil
}

func checkWinner(status *client.Status, txt *gui.Text) {
	if status.LastGameStatus == "win" {
		txt.SetBgColor(gui.Green)
		txt.SetText("Congratulation, you wiped your opponent off the board")
	} else {
		txt.SetBgColor(gui.Red)
		txt.SetText("I'm sorry, it's clearly not your day. You lost")
	}
}
