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
	_, initErr := app.client.InitGame("Ant0n", "You can run, but you can't hide", true)
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
	txt := gui.NewText(50, 15, "Let's start", nil)
	ui.Draw(txt)
	boardState.InitialStates(myBoard, enemyBoard, myShips)
	go func() error {
		appError := app.gameCourse(myBoard, enemyBoard, txt, ui)
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
		if gameStatus.GameStatus == "game_in_progress" {
			break
		}
	}
	return nil
}

func (app *App) gameCourse(myBoard, enemyBoard *gui.Board, txt *gui.Text, ui *gui.GUI) error {
	status, statusErr := app.client.Status()
	if statusErr != nil {
		log.Println(statusErr)
		return statusErr
	}
	for status.GameStatus == "game_in_progress" {
		for status.ShouldFire == true {
			if len(status.OppShots) != 0 {
				boardState.EnemyShoot(myBoard, status)
			}
			txt.SetText("It's your turn!")
			txt.SetText("Press on opponent's coordinates to shoot")
			var coords string
			for {
				coords = enemyBoard.Listen(context.TODO())
				x, y := ConvertCords(coords)
				if boardState.EnemyState[x][y] == "" {
					break
				}
			}
			txt.SetText(fmt.Sprintf("Coordinate: %s", coords))
			fire, fireErr := app.client.Fire(coords)
			if fireErr != nil {
				log.Println(fireErr)
				return fireErr
			}
			boardState.MarkMyShoot(txt, enemyBoard, fire, coords)
			status, _ = app.client.Status()
			time.Sleep(time.Second * 1)
		}
		time.Sleep(time.Second * 1)
		status, _ = app.client.Status()
	}
	if status.LastGameStatus == "win" {
		txt.SetText("Congratulation, you wiped your opponent off the board")
	} else {
		txt.SetText("I'm sorry, it's clearly not your day. You lost")
	}
	return nil
}
