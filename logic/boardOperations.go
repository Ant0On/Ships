package logic

import (
	"Ships/client"
	gui "github.com/grupawp/warships-gui/v2"
	"strconv"
)

type BoardState struct {
	MyState    [10][10]gui.State
	EnemyState [10][10]gui.State
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
	enemyBoard := gui.NewBoard(130, 3, nil)
	yourNick := gui.NewText(20, 1, description.Nick, nil)
	enemyNick := gui.NewText(150, 1, description.Opponent, nil)
	yourDesc := gui.NewText(5, 26, description.Desc, nil)
	enemyDesc := gui.NewText(135, 26, description.OppDesc, nil)
	ui.Draw(myBoard)
	ui.Draw(enemyBoard)
	ui.Draw(yourNick)
	ui.Draw(enemyNick)
	ui.Draw(yourDesc)
	ui.Draw(enemyDesc)
	return myBoard, enemyBoard, ui
}

func (boardState *BoardState) MarkMyShoot(txt *gui.Text, enemyBoard *gui.Board, fireResult string, coords string) {
	if fireResult == "hit" || fireResult == "sunk" {
		txt.SetText(fireResult)
		x, y := ConvertCords(coords)
		boardState.EnemyState[x][y] = gui.Hit
		enemyBoard.SetStates(boardState.EnemyState)
	} else {
		txt.SetText(fireResult)
		x, y := ConvertCords(coords)
		boardState.EnemyState[x][y] = gui.Miss
		enemyBoard.SetStates(boardState.EnemyState)
	}
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
