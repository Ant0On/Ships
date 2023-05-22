package logic

import (
	"Ships/client"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"
)

type BoardState struct {
	MyState      [10][10]gui.State
	EnemyState   [10][10]gui.State
	Ships        [10][10]gui.State
	Ui           *gui.GUI
	MyBoard      *gui.Board
	AccuracyText *gui.Text
	TimerTxt     *gui.Text
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

func CreateBoard(description *client.Description) (*gui.Board, *gui.Board) {
	boardState.Ui = gui.NewGUI(true)
	myBoard := gui.NewBoard(1, 3, nil)
	enemyBoard := gui.NewBoard(95, 3, nil)
	yourNick := gui.NewText(20, 1, description.Nick, nil)
	enemyNick := gui.NewText(95, 1, description.Opponent, nil)

	boardState.Ui.Draw(myBoard)
	boardState.Ui.Draw(enemyBoard)
	boardState.Ui.Draw(yourNick)
	boardState.Ui.Draw(enemyNick)

	handleDesc(description.Desc, description.OppDesc)

	return myBoard, enemyBoard
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

func handleDesc(myDesc, enemyDesc string) {
	wrapMyDesc := strings.Split(wordwrap.WrapString(myDesc, 40), "\n")
	wrapEnemyDesc := strings.Split(wordwrap.WrapString(enemyDesc, 40), "\n")

	for i, desc := range wrapMyDesc {
		boardState.Ui.Draw(gui.NewText(2, 26+i, desc, nil))
	}
	for i, desc := range wrapEnemyDesc {
		boardState.Ui.Draw(gui.NewText(97, 26+i, desc, nil))
	}
}
func (boardState *BoardState) countAccuracy(hits, totalShoots int) float64 {
	return float64(hits) / float64(totalShoots)
}
func makeShips() []string {
	ch := make(chan string, 20)
	exit := make(chan struct{})
	boardState.makeBoard()

	ctx, cancel := context.WithCancel(context.Background())

	coords := make([]string, 20)
	var square string
	var availableCoords, fCoords []string
	var forNowBorders []int
	newShip := []int{4, 7, 10, 12, 14, 16, 17, 18, 19}

	go func() {
		for i := 0; i < len(coords); i++ {

			switch {
			case i == 0:
				square = boardState.MyBoard.Listen(context.TODO())
			case slices.Contains(newShip, i):
				square, availableCoords, forNowBorders = boardState.createNewShip(forNowBorders, availableCoords)
			default:
				square = boardState.create(availableCoords)
			}

			options, borders := possibilities(square, availableCoords, forNowBorders)
			availableCoords = append(options)
			forNowBorders = append(borders)
			ch <- square
			time.Sleep(time.Millisecond * 300)
		}
		close(ch)
	}()
	go func() {
		for v := range ch {
			fCoords = append(fCoords, v)
		}
		close(exit)
		cancel()
	}()

	boardState.Ui.Start(ctx, nil)

	<-exit

	return fCoords

}

func possibilities(square string, finalCoords []string, finalBorder []int) ([]string, []int) {
	x, y := ConvertCords(square)
	boardState.MyState[x][y] = gui.Ship

	available, checkX, checkY := checkShip(x, y)

	for i := 0; i < len(available); i++ {
		for j := 0; j < len(available[0]); j++ {
			if boardState.MyState[checkX[j]][checkY[i]] != gui.Ship && available[i][j] &&
				boardState.MyState[checkX[j]][checkY[i]] != gui.Hit {
				boardState.MyState[checkX[j]][checkY[i]] = gui.Miss
				finalCoords = append(finalCoords, convertToString(checkX[j], checkY[i]))
			}
			finalBorder = append(finalBorder, checkX[j], checkY[i])
		}
	}
	boardState.MyBoard.SetStates(boardState.MyState)
	return finalCoords, finalBorder
}

func checkShip(x, y int) ([][]bool, []int, []int) {

	available := availableChecker(x, y)
	checkX := xChecker(x)
	checkY := yChecker(y)
	return available, checkX, checkY
}
