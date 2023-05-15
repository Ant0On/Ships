package logic

import (
	"Ships/client"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/mitchellh/go-wordwrap"
	"strconv"
	"strings"
	"time"
)

type BoardState struct {
	MyState      [10][10]gui.State
	EnemyState   [10][10]gui.State
	AccuracyText *gui.Text
	Accuracy     float64
}

func ConvertCords(cords string) (int, int) {
	x := int(cords[0] - 'A')
	y, _ := strconv.Atoi(cords[1:])
	y--
	return x, y
}
func PlaceShips(cords []string, states [][]gui.State) {
	for _, cord := range cords {
		x, y := ConvertCords(cord)
		states[x][y] = gui.Ship
	}
}

func CreateBoard(description *client.Description) (*gui.Board, *gui.Board, *gui.GUI) {
	ui := gui.NewGUI(true)
	myBoard := gui.NewBoard(1, 3, nil)
	enemyBoard := gui.NewBoard(95, 3, nil)
	yourNick := gui.NewText(20, 1, description.Nick, nil)
	enemyNick := gui.NewText(95, 1, description.Opponent, nil)

	ui.Draw(myBoard)
	ui.Draw(enemyBoard)
	ui.Draw(yourNick)
	ui.Draw(enemyNick)

	handleDesc(ui, description.Desc, description.OppDesc)

	return myBoard, enemyBoard, ui
}

func (boardState *BoardState) MarkMyShoot(mark *gui.Text, enemyBoard *gui.Board, fireResult string, coords string) {
	if fireResult == "hit" || fireResult == "sunk" {
		mark.SetBgColor(gui.Red)
		mark.SetText(fireResult)
		x, y := ConvertCords(coords)
		boardState.EnemyState[x][y] = gui.Hit
		enemyBoard.SetStates(boardState.EnemyState)
	} else {
		mark.SetBgColor(gui.Blue)
		mark.SetText(fireResult)
		x, y := ConvertCords(coords)
		boardState.EnemyState[x][y] = gui.Miss
		enemyBoard.SetStates(boardState.EnemyState)
	}
	boardState.AccuracyText.SetBgColor(gui.NewColor(184, 27, 227))
	boardState.AccuracyText.SetText(fmt.Sprintf("AccuracyText: %.2f", boardState.Accuracy))
	time.Sleep(time.Second * 2)

}
func (boardState *BoardState) InitialStates(myBoard, enemyBoard *gui.Board, coords []string) {

	for i := range boardState.MyState {
		boardState.MyState[i] = [10]gui.State{}
		boardState.EnemyState[i] = [10]gui.State{}
	}
	stateSlice := make([][]gui.State, len(boardState.MyState))
	for i := range stateSlice {
		stateSlice[i] = boardState.MyState[i][:]
	}
	PlaceShips(coords, stateSlice)
	myBoard.SetStates(boardState.MyState)
	enemyBoard.SetStates(boardState.EnemyState)

}
func (boardState *BoardState) EnemyShoot(myBoard *gui.Board, status *client.Status) {
	for _, coords := range status.OppShots {
		x, y := ConvertCords(coords)

		switch state := &boardState.MyState[x][y]; *state {
		case gui.Hit, gui.Ship:
			*state = gui.Hit
		default:
			*state = gui.Miss
		}
	}
	myBoard.SetStates(boardState.MyState)
}

func handleDesc(ui *gui.GUI, myDesc, enemyDesc string) {
	wrapMyDesc := strings.Split(wordwrap.WrapString(myDesc, 40), "\n")
	wrapEnemyDesc := strings.Split(wordwrap.WrapString(enemyDesc, 40), "\n")

	for i, desc := range wrapMyDesc {
		ui.Draw(gui.NewText(2, 26+i, desc, nil))
	}
	for i, desc := range wrapEnemyDesc {
		ui.Draw(gui.NewText(97, 26+i, desc, nil))
	}
}
func (boardState *BoardState) countAccuracy(hits, totalShoots int) float64 {
	return float64(hits) / float64(totalShoots)
}
