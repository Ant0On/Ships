package logic

import (
	"context"
	gui "github.com/grupawp/warships-gui/v2"
	"golang.org/x/exp/slices"
	"strconv"
)

func boardConf() *gui.BoardConfig {
	boardConfig := gui.NewBoardConfig()
	boardConfig.MissColor = gui.NewColor(42, 145, 19)
	boardConfig.MissChar = 'A'
	boardConfig.HitChar = 'N'
	return boardConfig
}

func convertToString(x, y int) string {
	xChar := byte('A' + x)
	yChar := strconv.Itoa(y + 1)
	return string(xChar) + yChar
}

func (boardState *BoardState) makeBoard() {
	boardState.Ui = gui.NewGUI(true)
	boardConfig := boardConf()
	boardState.MyBoard = gui.NewBoard(10, 3, boardConfig)
	boardState.Ui.Draw(boardState.MyBoard)
	txt := gui.NewText(50, 5, "", nil)
	boardState.Ui.Draw(txt)
}

func (boardState *BoardState) create(availableCoords []string) string {
	var square string
	for {
		square = boardState.MyBoard.Listen(context.TODO())
		if slices.Contains(availableCoords, square) {
			break
		}
	}
	return square
}
func (boardState *BoardState) createNewShip(borders []int, availableCoords []string) (string, []string, []int) {

	var square string
	for k := 0; k < len(borders); k += 2 {
		if boardState.Ships[borders[k]][borders[k+1]] != gui.Ship {
			boardState.Ships[borders[k]][borders[k+1]] = gui.Hit
		}
	}
	boardState.MyBoard.SetStates(boardState.Ships)
	for {
		square = boardState.MyBoard.Listen(context.TODO())
		x, y := ConvertCords(square)
		if boardState.Ships[x][y] != gui.Ship && boardState.Ships[x][y] != gui.Hit {
			break
		}
	}
	availableCoords = []string{}
	borders = []int{}
	return square, availableCoords, borders
}

func xChecker(x int) []int {
	switch {
	case x == 0:
		return []int{x, x + 1}
	case x == 9:
		return []int{x - 1, x}
	default:
		return []int{x - 1, x, x + 1}
	}
}

func yChecker(y int) []int {
	switch {
	case y == 0:
		return []int{y, y + 1}
	case y == 9:
		return []int{y - 1, y}
	default:
		return []int{y - 1, y, y + 1}
	}
}

func availableChecker(x, y int) [][]bool {
	switch {
	case upperLeft(x, y):
		return [][]bool{
			{false, true},
			{true, false},
		}
	case bottomLeft(x, y):
		return [][]bool{
			{true, false},
			{false, true},
		}
	case upperBorder(x, y):
		return [][]bool{
			{true, false, true},
			{false, true, false},
		}
	case lowerBorder(x, y):
		return [][]bool{
			{false, true, false},
			{true, false, true},
		}
	case leftBorder(x, y):
		return [][]bool{
			{true, false},
			{false, true},
			{true, false},
		}
	case rightBorder(x, y):
		return [][]bool{
			{false, true},
			{true, false},
			{false, true},
		}
	default:
		return [][]bool{
			{false, true, false},
			{true, false, true},
			{false, true, false},
		}
	}
}

func upperLeft(x, y int) bool {
	border := []int{0, 9}
	if x == y && slices.Contains(border, x) {
		return true
	}
	return false
}
func bottomLeft(x, y int) bool {
	border := []int{0, 9}
	if x != y && slices.Contains(border, x) && slices.Contains(border, y) {
		return true
	}
	return false
}
func upperBorder(x, y int) bool {
	if x > 0 && y == 0 {
		return true
	}
	return false
}
func lowerBorder(x, y int) bool {
	if x > 0 && y == 9 {
		return true
	}
	return false
}
func leftBorder(x, y int) bool {
	if x == 0 && y > 0 {
		return true
	}
	return false
}
func rightBorder(x, y int) bool {
	if x == 9 && y > 0 {
		return true
	}
	return false
}
