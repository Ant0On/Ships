package logic

import (
	"Ships/client"
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"time"
)

func (app *App) generateMenu() error {
	pterm.Println() //
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString("Battleships")).Srender()
	pterm.DefaultCenter.Println(s)
	pterm.Println() //
	for {
		endGame := true
		options := []string{
			"Show potential enemies",
			"Start game",
			"Show top 10 players",
			"Show player stats",
			"Quit",
		}
		pterm.Println() //
		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
		pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))
		pterm.Println() //
		switch selectedOption {
		case options[0]:
			time.Sleep(time.Second * 1)
			pterm.Info.Println("Potential enemies: ")
			playerListErr := app.client.PlayersList()
			if playerListErr != nil {
				return playerListErr
			}
			pterm.Println() //
		case options[1]:
			var coords []string
			info, bot := initialConfig()
			pterm.Info.Print("Do you want to set your ships?: ")
			result, _ := pterm.DefaultInteractiveConfirm.Show()
			pterm.Println() // Blank line
			pterm.Info.Printfln("You answered: %s", boolToText(result))
			if result == true {
				coords = makeShips()
			}
			initErr := app.client.InitGame(client.InitialData{Coords: coords, Desc: info[1],
				Nick: info[0], TargetNick: info[2], Wpbot: bot})
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
			time.Sleep(time.Millisecond * 300)
			createBoard(description)
			boardState.initialStates(myShips)
			appError := app.gameCourse()
			if appError != nil {
				log.Println(appError)
				return appError
			}
			fmt.Print("Do you want to start another game? ")
			endGame, _ = pterm.DefaultInteractiveConfirm.Show()
			pterm.Println() // Blank line
			pterm.Info.Printfln("You answered: %s", boolToText(endGame))
			pterm.Println() //
			if endGame == false {
				break
			}
			time.Sleep(time.Second * 3)
		case options[2]:
			pterm.Info.Println("Top 10 players: ")
			playerListErr := app.client.Top10()
			if playerListErr != nil {
				return playerListErr
			}
			pterm.Println() //
		case options[3]:
			pterm.Info.Println("Enter the player's nickname: ")
			nickName, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
			pterm.Println() //
			playerListErr := app.client.PlayerStats(nickName)
			pterm.Println() //
			if playerListErr != nil {
				return playerListErr
			}
			pterm.Println() //
		case options[4]:
			endGame = false
		}
		if endGame == false {
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
