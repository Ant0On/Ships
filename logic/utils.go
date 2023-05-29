package logic

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/mitchellh/go-wordwrap"
	"github.com/pterm/pterm"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
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
		if boardState.MyState[borders[k]][borders[k+1]] != gui.Ship {
			boardState.MyState[borders[k]][borders[k+1]] = gui.Hit
		}
	}
	boardState.MyBoard.SetStates(boardState.MyState)
	for {
		square = boardState.MyBoard.Listen(context.TODO())
		x, y := convertCords(square)
		if boardState.MyState[x][y] != gui.Ship && boardState.MyState[x][y] != gui.Hit {
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

func handleDesc(myDesc, enemyDesc string) {
	wrapMyDesc := strings.Split(wordwrap.WrapString(myDesc, 40), "\n")
	wrapEnemyDesc := strings.Split(wordwrap.WrapString(enemyDesc, 40), "\n")

	for i, desc := range wrapMyDesc {
		boardState.Ui.Draw(gui.NewText(7, 33+i, desc, nil))
	}
	for i, desc := range wrapEnemyDesc {
		boardState.Ui.Draw(gui.NewText(100, 33+i, desc, nil))
	}
}

func countAccuracy(hits, totalShoots int) float64 {
	return float64(hits) / float64(totalShoots)
}

func convertCords(cords string) (int, int) {
	x := int(cords[0] - 'A')
	y, _ := strconv.Atoi(cords[1:])
	y--
	return x, y
}

func boardPrepare() {
	gameConf.myTour = gui.NewText(18, 6, "", nil)
	gameConf.enemyTour = gui.NewText(99, 6, "", nil)
	gameConf.generalInfo = gui.NewText(50, 15, "Let's start", nil)
	gameConf.timer = gui.NewText(65, 10, "", nil)
	gameConf.countFour = 1
	gameConf.countThree = 2
	gameConf.countTwo = 3
	gameConf.countSingle = 4
	gameConf.fourMasted = gui.NewText(50, 17, fmt.Sprintf("Four-masted left: %d", gameConf.countFour), nil)
	gameConf.threeMasted = gui.NewText(50, 19, fmt.Sprintf("Three-masted left: %d", gameConf.countThree), nil)
	gameConf.twoMasted = gui.NewText(50, 21, fmt.Sprintf("Two-masted left: %d", gameConf.countTwo), nil)
	gameConf.singleMasted = gui.NewText(50, 23, fmt.Sprintf("Single-masted left: %d", gameConf.countSingle), nil)
	exitConf := gui.NewTextConfig()
	exitConf.BgColor = gui.NewColor(117, 11, 11)
	exitConf.FgColor = gui.White
	exitInfo := gui.NewText(55, 30, "If you want to exit press ctrl+c", exitConf)
	boardState.Ui.Draw(gameConf.myTour)
	boardState.Ui.Draw(gameConf.enemyTour)
	boardState.Ui.Draw(gameConf.generalInfo)
	boardState.Ui.Draw(gameConf.timer)
	boardState.Ui.Draw(exitInfo)
	boardState.Ui.Draw(gameConf.fourMasted)
	boardState.Ui.Draw(gameConf.threeMasted)
	boardState.Ui.Draw(gameConf.twoMasted)
	boardState.Ui.Draw(gameConf.singleMasted)
}

func legend() {
	conf := gui.NewTextConfig()
	conf.BgColor = gui.NewColor(245, 242, 242)

	generateBorder()
	conf.BgColor = gui.Red
	boardState.Ui.Draw(gui.NewText(150, 13, "Red square with letter H - HIT", conf))
	conf.BgColor = gui.Grey
	boardState.Ui.Draw(gui.NewText(150, 15, "Grey square with letter M - MISS", conf))
	conf.BgColor = gui.NewColor(52, 227, 218)
	boardState.Ui.Draw(gui.NewText(150, 17, "On enemy board", conf))
	conf.BgColor = gui.NewColor(224, 208, 31)
	boardState.Ui.Draw(gui.NewText(150, 19, "Yellow square with letter S - SUNK", conf))
	conf.BgColor = gui.NewColor(52, 227, 218)
	boardState.Ui.Draw(gui.NewText(150, 21, "On my board", conf))
	conf.BgColor = gui.NewColor(224, 208, 31)
	boardState.Ui.Draw(gui.NewText(150, 23, "Yellow square with letter S - SHIP", conf))

}

func generateBorder() {
	conf := gui.NewTextConfig()
	conf.BgColor = gui.NewColor(245, 242, 242)
	boardState.Ui.Draw(gui.NewText(170, 9, "Legend", conf))
	boardState.Ui.Draw(gui.NewText(149, 11, "-------------------------------------------------", conf))
	for i := 11; i <= 25; i++ {
		boardState.Ui.Draw(gui.NewText(148, i, "|", conf))
		boardState.Ui.Draw(gui.NewText(198, i, "|", conf))
	}
	boardState.Ui.Draw(gui.NewText(149, 25, "-------------------------------------------------", conf))
}

func informAboutShip(i int, txt *gui.Text) {
	switch i {
	case 4:
		txt.SetText("Create first three-masted ship")
	case 7:
		txt.SetText("Create a second three-masted ship")
	case 10:
		txt.SetText("Create first two-masted ship")
	case 12:
		txt.SetText("Create a second two-masted ship")
	case 14:
		txt.SetText("Create the third two-masted ship")
	case 16:
		txt.SetText("Create first single-masted ship")
	case 17:
		txt.SetText("Create a second single-masted ship")
	case 18:
		txt.SetText("Create the third single-masted ship")
	case 19:
		txt.SetText("Create the last single-masted ship")
	}
}
func WrapEndInfo(txt string, config *gui.TextConfig) {
	wrapTxt := strings.Split(wordwrap.WrapString(txt, 40), "\n")

	for i, ch := range wrapTxt {
		boardState.Ui.Draw(gui.NewText(50, 10+i, ch, config))
	}
}
func validateName(nickName string) string {
	for len(nickName) != 1 || len(nickName) > 10 {
		pterm.Info.Println("Nick should be in range (2 -10) ")
		pterm.Info.Println("Enter your nickname: ")
		nickName, _ = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	}
	return nickName
}

func checkNeighbour(x, y int) int {
	totalCount := 0
	available, checkX, checkY := checkShip(x, y)
	for i := 0; i < len(available); i++ {
		for j := 0; j < len(available[0]); j++ {
			if boardState.EnemyState[checkX[j]][checkY[i]] == gui.Hit {
				boardState.EnemyState[checkX[j]][checkY[i]] = gui.Ship
				totalCount++
				return totalCount + checkNeighbour(checkX[j], checkY[i])
			}
		}
	}
	return totalCount
}

func markSunk(neighbours int) {
	switch neighbours {
	case 3:
		gameConf.countFour--
		gameConf.fourMasted.SetText(fmt.Sprintf("Four-masted left: %d", gameConf.countFour))
	case 2:
		gameConf.countThree--
		gameConf.threeMasted.SetText(fmt.Sprintf("Three-masted left: %d", gameConf.countThree))
	case 1:
		gameConf.countTwo--
		gameConf.twoMasted.SetText(fmt.Sprintf("Two-masted left: %d", gameConf.countTwo))
	case 0:
		gameConf.countSingle--
		gameConf.singleMasted.SetText(fmt.Sprintf("Single-masted left: %d", gameConf.countSingle))
	}
}
