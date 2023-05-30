package logic

import (
	"Ships/client"
	"context"
	gui "github.com/grupawp/warships-gui/v2"
	"golang.org/x/exp/slices"
	"time"
)

type BoardState struct {
	Ui         *gui.GUI
	MyBoard    *gui.Board
	EnemyState [10][10]gui.State
	MyState    [10][10]gui.State
	EnemyBoard *gui.Board
}

const (
	missState = gui.Miss
	hitState  = gui.Hit
	shipState = gui.Ship
)

func placeShips(cords []string, states [][]gui.State) {
	for _, cord := range cords {
		x, y := convertCords(cord)
		states[x][y] = gui.Ship
	}
}

func createBoard(description *client.Description) {
	boardState.Ui = gui.NewGUI(true)
	boardState.MyBoard = gui.NewBoard(3, 10, nil)
	boardState.EnemyBoard = gui.NewBoard(98, 10, nil)
	yourNick := gui.NewText(20, 3, description.Nick, nil)
	enemyNick := gui.NewText(110, 3, description.Opponent, nil)

	boardState.Ui.Draw(boardState.MyBoard)
	boardState.Ui.Draw(boardState.EnemyBoard)
	boardState.Ui.Draw(yourNick)
	boardState.Ui.Draw(enemyNick)

	handleDesc(description.Desc, description.OppDesc)

}

func (boardState *BoardState) markMyShoot(fireResult string, coords string) {
	if fireResult == "hit" {
		markShootConf(hitState, coords)
	} else if fireResult == "sunk" {
		x, y := convertCords(coords)
		neighbours := checkNeighbour(x, y)
		markSunk(neighbours)
		markShootConf(shipState, coords)
	} else {
		markShootConf(missState, coords)
	}

}
func (boardState *BoardState) initialStates(coords []string) {
	boardState.MyState = [10][10]gui.State{}
	boardState.EnemyState = [10][10]gui.State{}

	for i := range boardState.MyState {
		boardState.MyState[i] = [10]gui.State{}
		boardState.EnemyState[i] = [10]gui.State{}
	}
	stateSlice := make([][]gui.State, len(boardState.MyState))
	for i := range stateSlice {
		stateSlice[i] = boardState.MyState[i][:]
	}
	placeShips(coords, stateSlice)
	boardState.MyBoard.SetStates(boardState.MyState)
	boardState.EnemyBoard.SetStates(boardState.EnemyState)

}
func (boardState *BoardState) enemyShoot(status *client.Status) {
	for _, coords := range status.OppShots {
		x, y := convertCords(coords)

		switch state := &boardState.MyState[x][y]; *state {
		case hitState, shipState:
			*state = hitState
		default:
			*state = missState
		}
	}
	boardState.MyBoard.SetStates(boardState.MyState)
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
	legendInfo := gui.NewText(60, 5, "Create a four-masted ship", nil)
	boardState.Ui.Draw(legendInfo)

	go func() {
		for i := 0; i < len(coords); i++ {

			switch {
			case i == 0:
				square = boardState.MyBoard.Listen(context.TODO())
			case slices.Contains(newShip, i):
				informAboutShip(i, legendInfo)
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
	x, y := convertCords(square)
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

func markShootConf(state gui.State, fireInfo string) {
	x, y := convertCords(fireInfo)
	boardState.EnemyState[x][y] = state
	boardState.EnemyBoard.SetStates(boardState.EnemyState)

}
