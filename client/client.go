package client

import (
	"Ships/models"
	"bytes"
	"encoding/json"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseURL string, timeOut time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeOut,
		},
	}
}
func (client *Client) InitGame() error {
	initialData := &models.InitialData{Wpbot: true}

	body, bodyErr := json.Marshal(initialData)
	if bodyErr != nil {
		log.Println(bodyErr)
	}

	initUrl, urlErr := url.JoinPath(client.BaseURL, "/game")
	if urlErr != nil {
		log.Println(urlErr)
	}

	req, reqErr := http.NewRequest("POST", initUrl, bytes.NewReader(body))
	if reqErr != nil {
		log.Println(reqErr)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, respErr := client.HTTPClient.Do(req)
	if respErr != nil {
		log.Println(respErr)
	}
	client.Token = resp.Header.Get("X-Auth-Token")

	return reqErr
}
func (client *Client) Board() ([]string, error) {
	board := models.Board{}
	boardURL, urlErr := url.JoinPath(client.BaseURL, "/game/board")
	fmt.Println(boardURL)
	if urlErr != nil {
		log.Println(urlErr)
	}

	req, reqErr := http.NewRequest("GET", boardURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", client.Token)
	if reqErr != nil {
		log.Println(reqErr)
	}

	res, resErr := client.HTTPClient.Do(req)
	if resErr != nil {
		log.Println(resErr)
	}
	defer res.Body.Close()

	cords, cordsErr := io.ReadAll(res.Body)
	if cordsErr != nil {
		log.Println(cordsErr)
	}
	fmt.Println(string(cords))

	jsonErr := json.Unmarshal(cords, &board)
	fmt.Println(board.Board)
	return board.Board, jsonErr

}

func convertCords(cords string) (int, int) {
	x := int(cords[0] - 'A')
	y, _ := strconv.Atoi(cords[1:])
	y--
	return x, y
}

func PlaceShips(cords []string, states [][]gui.State) {
	for _, cord := range cords {
		x, y := convertCords(cord)
		states[x][y] = gui.Ship
	}
}

func (client *Client) Status() (*http.Response, *models.Status, error) {
	stats := models.Status{}
	statusUrl, urlErr := url.JoinPath(client.BaseURL, "/game")
	if urlErr != nil {
		log.Println(urlErr)
	}
	req, reqErr := http.NewRequest("GET", statusUrl, nil)
	req.Header.Set("X-Auth-Token", client.Token)
	if reqErr != nil {
		log.Println(reqErr)
	}
	res, resErr := client.HTTPClient.Do(req)
	if resErr != nil {
		log.Println(resErr)
	}
	defer res.Body.Close()

	data, dataErr := io.ReadAll(res.Body)
	if dataErr != nil {
		log.Println(dataErr)
	}
	json.Unmarshal(data, &stats)
	return res, &stats, resErr
}

func (client *Client) Shoot(cord string, status *models.Status) (string, error) {
	fireData := models.Fire{Coord: cord}
	body, bodyErr := json.Marshal(fireData)
	if bodyErr != nil {
		log.Println(bodyErr)
	}
	if status.ShouldFire == true {
		fmt.Println("Your turn:")
		var w1 string
		n, err := fmt.Scanln(&w1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(n)

	}
	fireUrl, fireErr := url.JoinPath(client.BaseURL, "/game/fire")
	if fireErr != nil {
		log.Println(fireErr)
	}
	req, reqErr := http.NewRequest("POST", fireUrl, bytes.NewReader(body))
	req.Header.Set("X-Auth-Token", client.Token)
	if reqErr != nil {
		log.Println(reqErr)
	}
	return "", nil
}
