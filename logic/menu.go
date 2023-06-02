package logic

import (
	"Ships/client"
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"time"
)

const (
	potentialEnemies = "Show potential enemies"
	startGame        = "Start game"
	top10            = "Show top 10 players"
	playerStats      = "Show player stats"
	quit             = "Quit"
)

func (app *App) generateMenu() error {
	pterm.Println() //
	s, err := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString("Battleships")).Srender()
	if err != nil {
		return err
	}
	pterm.DefaultCenter.Println(s)
	pterm.Println() //
	for {
		endGame := true
		options := []string{
			potentialEnemies,
			startGame,
			top10,
			playerStats,
			quit,
		}
		pterm.Println() //
		selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
		if err != nil {
			return err
		}
		pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))
		pterm.Println() //
		switch selectedOption {
		case potentialEnemies:
			time.Sleep(time.Second * 1)
			pterm.Info.Println("Potential enemies: ")
			err := app.client.PlayersList()
			if err != nil {
				return err
			}
			pterm.Println()
		case startGame:
			var coords []string
			info, bot := initialConfig()
			pterm.Info.Print("Do you want to set your ships?: ")
			result, err := pterm.DefaultInteractiveConfirm.Show()
			if err != nil {
				return err
			}
			pterm.Println() // Blank line
			pterm.Info.Printfln("You answered: %s", boolToText(result))
			if result == true {
				coords = makeShips()
				if len(coords) != 20 {
					coords = []string{}
				}
			}
			time.Sleep(time.Second * 1)
			err = app.client.InitGame(client.InitialData{Coords: coords, Desc: info[1],
				Nick: info[0], TargetNick: info[2], Wpbot: bot})
			if err != nil {
				log.Println(err)
				return err
			}
			time.Sleep(time.Second * 1)
			err = app.waitToStart()
			if err != nil {
				log.Println(err)
				return err
			}
			time.Sleep(time.Second * 1)
			description, err := app.client.Descriptions()
			if err != nil {
				log.Println(err)
				return err
			}
			time.Sleep(time.Second * 1)
			myShips, err := app.client.Board()
			if err != nil {
				log.Println(err)
				return err
			}
			time.Sleep(time.Millisecond * 300)
			createBoard(description)
			boardState.initialStates(myShips)
			err = app.gameCourse()
			if err != nil {
				log.Println(err)
				return err
			}
			fmt.Print("Do you want to start another game? ")
			endGame, err = pterm.DefaultInteractiveConfirm.Show()
			if err != nil {
				return err
			}
			pterm.Println() // Blank line
			pterm.Info.Printfln("You answered: %s", boolToText(endGame))
			pterm.Println() //
			if !endGame {
				break
			}
			time.Sleep(time.Second * 3)
		case top10:
			pterm.Info.Println("Top 10 players: ")
			err := app.client.Top10()
			if err != nil {
				return err
			}
			pterm.Println() //
		case playerStats:
			pterm.Info.Println("Enter the player's nickname: ")
			nickName, err := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
			if err != nil {
				return err
			}
			pterm.Println() //
			err = app.client.PlayerStats(nickName)
			pterm.Println() //
			if err != nil {
				return err
			}
			pterm.Println() //
		case quit:
			endGame = false
		}
		if !endGame {
			break
		}
	}
	return nil
}
func boolToText(b bool) string {
	if b {
		return pterm.Green("Yes")
	}
	return pterm.Red("No")
}
