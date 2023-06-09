package logic

import (
	"Ships/client"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/pterm/pterm"
	"log"
	"strconv"
	"time"
)

type GameConf struct {
	myTour       *gui.Text
	enemyTour    *gui.Text
	generalInfo  *gui.Text
	timer        *gui.Text
	mark         *gui.Text
	fourMasted   *gui.Text
	threeMasted  *gui.Text
	twoMasted    *gui.Text
	singleMasted *gui.Text
	countFour    int
	countThree   int
	countTwo     int
	countSingle  int
}

type App struct {
	client *client.Client
}

var boardState BoardState

var gameConf GameConf

func NewApp(client *client.Client) *App {
	return &App{client: client}
}

func (app *App) Run() error {
	err := app.generateMenu()
	if err != nil {
		return err
	}
	return nil
}

func (app *App) waitToStart() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go app.everyThreeSecond(ctx)
	for {
		time.Sleep(time.Second * 1)
		gameStatus, err := app.client.Status()
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println(gameStatus.GameStatus)

		if gameStatus.GameStatus == "game_in_progress" {
			break
		}
	}
	return nil
}
func (app *App) everyThreeSecond(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := app.client.Refresh()
			if err != nil {
				return
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
func (app *App) gameCourse() error {
	exit := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	status, err := app.client.Status()
	if err != nil {
		log.Fatal(err)
	}

	fireRes := ""
	hits := 0
	totalShoots := 0
	boardPrepare()
	legend()

	gameConf.mark = gui.NewText(65, 5, "", nil)
	boardState.Ui.Draw(gameConf.mark)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			timerStatus, _ := app.client.Status()
			for i := 60; timerStatus.ShouldFire == true; i-- {
				gameConf.timer.SetText(fmt.Sprintf("Time left: %d", i))
				time.Sleep(time.Millisecond * 900)
				timerStatus, _ = app.client.Status()
			}
			boardState.enemyShoot(status)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			for status, err := app.client.Status(); status.ShouldFire == true; status, err = app.client.Status() {
				if err != nil {
					log.Fatal(err)
				}
				gameConf.enemyTour.SetText("")
				boardState.enemyShoot(status)
				gameConf.myTour.SetBgColor(gui.White)
				gameConf.myTour.SetText("Your turn!")
				gameConf.generalInfo.SetText("Press on opponent's coordinates to shoot")
				var coords string
				for {
					coords = boardState.EnemyBoard.Listen(context.TODO())
					x, y := convertCords(coords)
					if boardState.EnemyState[x][y] == "" {
						break
					}
				}
				time.Sleep(time.Second * 1)
				fire, err := app.client.Fire(coords)
				if err != nil {
					log.Println(err)
					return
				}
				gameConf.mark.SetBgColor(gui.White)
				gameConf.mark.SetText(fmt.Sprintf("Coordinate: %s", coords))
				boardState.markMyShoot(fire, coords)
				fireRes = fire
			}
			gameConf.myTour.SetText("")
			gameConf.enemyTour.SetText("Opponent's turn")
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			time.Sleep(time.Millisecond * 100)

			if fireRes != "miss" && fireRes != "" {
				hits++
			}
			if fireRes != "" {
				totalShoots++
			}
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			time.Sleep(time.Millisecond * 100)
			status, err := app.client.Status()
			if err != nil {
				log.Fatal(err)
			}
			if status.GameStatus == "ended" {
				boardState.enemyShoot(status)
				accuracy := countAccuracy(hits, totalShoots)
				str := strconv.FormatFloat(accuracy*100, 'f', 2, 64)
				str += "% accuracy"
				accuracyTxt := gui.NewText(65, 13, str, nil)
				boardState.Ui.Draw(accuracyTxt)
				checkWinner(status)
				time.Sleep(time.Second * 10)
				close(exit)
				cancel()
				break
			}
		}
	}()

	boardState.Ui.Start(ctx, nil)
	select {
	case value := <-exit:
		fmt.Println(value)
		app.client.HTTPClient.CloseIdleConnections()
		<-exit
	default:
		close(exit)
		err := app.client.Abandon()
		if err != nil {
			return err
		}
	}

	return nil
}

func checkWinner(status *client.Status) bool {
	win := "Congratulation, you wiped your opponent off the board"
	lose := "I'm sorry, it's clearly not your day. You lost"
	config := gui.NewTextConfig()

	if status.LastGameStatus == "win" {
		config.BgColor = gui.NewColor(25, 196, 22)
		config.FgColor = gui.Black
		WrapEndInfo(win, config)
		return true
	} else {
		config.BgColor = gui.NewColor(117, 11, 11)
		config.FgColor = gui.White
		gameConf.generalInfo.SetBgColor(gui.Red)
		WrapEndInfo(lose, config)
		return false
	}
}

func initialConfig() ([]string, bool) {
	var nickName, targetNick, description string
	var decision bool
	pterm.Info.Println("Enter your nickname: ")
	nickName, _ = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	nickName = validateName(nickName)
	pterm.Println()
	pterm.Info.Println("Enter your description: ")
	description, _ = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	pterm.Println()
	pterm.Info.Println("Do you want to fight against bot? ")
	decision, _ = pterm.DefaultInteractiveConfirm.Show()
	pterm.Println()
	if decision {
		return []string{nickName, description, ""}, decision
	}
	pterm.Info.Println("Enter your opponent's nickname: ")
	targetNick, _ = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	return []string{nickName, description, targetNick}, decision
}
