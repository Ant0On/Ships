package logic

import (
	"Ships/client"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/pterm/pterm"
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
	generateErr := app.generateMenu()
	if generateErr != nil {
		return generateErr
	}
	return nil
}

func (app *App) waitToStart() error {
	stop := make(chan bool)
	go app.everyThreeSecond(stop)
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
func (app *App) everyThreeSecond(stop chan bool) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			refreshErr := app.client.Refresh()
			if refreshErr != nil {
				return
			}
		case <-stop:
			ticker.Stop()
			return
		}
	}
}
func (app *App) gameCourse(myBoard, enemyBoard *gui.Board) error {
	ch := make(chan bool)
	exit := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	status, statusErr := app.client.Status()
	if statusErr != nil {
		log.Println(statusErr)
		return statusErr
	}
	fireRes := ""
	hits := 0
	totalShoots := 0
	txt := gui.NewText(50, 5, "Let's start", nil)
	boardState.AccuracyText = gui.NewText(5, 30, "Accuracy: 0", nil)
	timer := gui.NewText(50, 20, "", nil)
	boardState.Ui.Draw(timer)
	boardState.Ui.Draw(txt)
	mark := gui.NewText(65, 15, "", nil)
	boardState.Ui.Draw(mark)
	time.Sleep(time.Millisecond * 200)

	go func() {
		for {
			timerStatus, _ := app.client.Status()
			time.Sleep(time.Millisecond / 4)
			if timerStatus.ShouldFire == true {
				timer.SetText(fmt.Sprintf("Time left: %d", timerStatus.Timer))
			} else {
				timer.SetText("")
			}
			boardState.EnemyShoot(myBoard, status)
		}
	}()

	go func() {
		for {
			status, _ = app.client.Status()
			for status.ShouldFire == true {
				boardState.EnemyShoot(myBoard, status)
				txt.SetBgColor(gui.White)
				txt.SetText("It's your turn!")
				time.Sleep(time.Millisecond * 200)
				txt.SetText("Press on opponent's coordinates to shoot")
				var coords string
				for {
					coords = enemyBoard.Listen(context.TODO())
					x, y := ConvertCords(coords)
					if boardState.EnemyState[x][y] == "" {
						break
					}
				}
				fire, fireErr := app.client.Fire(coords)
				if fireErr != nil {
					log.Println(fireErr)
					return
				}
				mark.SetBgColor(gui.White)
				mark.SetText(fmt.Sprintf("Coordinate: %s", coords))
				boardState.MarkMyShoot(mark, enemyBoard, fire, coords)
				fireRes = fire
				status, _ = app.client.Status()
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 200)

			if fireRes != "miss" && fireRes != "" {
				hits++
			}
			if fireRes != "" {
				totalShoots++
			}
			boardState.Accuracy = boardState.countAccuracy(hits, totalShoots)
		}
	}()
	go func() {
		for {
			if status.GameStatus == "ended" {
				boardState.EnemyShoot(myBoard, status)
				boardState.Ui.Draw(boardState.AccuracyText)
				checkWinner(status, txt)
				time.Sleep(time.Second * 2)
				ch <- checkWinner(status, txt)
				close(ch)
				break
			}
		}
	}()
	go func() {
		for v := range ch {
			log.Println(v)
		}
		close(exit)
		cancel()
	}()
	boardState.Ui.Start(ctx, nil)
	app.client.HTTPClient.CloseIdleConnections()
	<-exit
	return nil
}

func checkWinner(status *client.Status, txt *gui.Text) bool {
	if status.LastGameStatus == "win" {
		txt.SetBgColor(gui.Green)
		txt.SetText("Congratulation, you wiped your opponent off the board")
		return true
	} else {
		txt.SetBgColor(gui.Red)
		txt.SetText("I'm sorry, it's clearly not your day. You lost")
		return false
	}
}

func initialConfig() ([]string, bool) {
	var nickName, targetNick, description string
	var decision bool
	pterm.Info.Println("Enter your nickname: ")
	nickName, _ = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	pterm.Println() // Blank line
	pterm.Info.Println("Enter your description: ")
	description, _ = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	pterm.Println() // Blank line
	pterm.Info.Println("Do you want to fight against bot? ")
	decision, _ = pterm.DefaultInteractiveConfirm.Show()
	pterm.Println() // Blank line
	if decision == true {
		return []string{nickName, description, ""}, decision
	}
	pterm.Info.Println("Enter your opponent's nickname: ")
	targetNick, _ = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	return []string{nickName, description, targetNick}, decision
}
