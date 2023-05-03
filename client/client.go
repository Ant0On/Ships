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
func (client *Client) InitGame(desc, nick string, wpBot bool) (*models.InitialData, error) {
	initialData := &models.InitialData{
		Desc:  desc,
		Nick:  nick,
		Wpbot: wpBot,
	}

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

	return initialData, respErr
}
func (client *Client) Board() ([]string, error) {
	board := models.Board{}
	boardURL, urlErr := url.JoinPath(client.BaseURL, "/game/board")

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

func (client *Client) Status() (*models.Status, error) {
	stats := models.Status{}
	statusUrl, urlErr := url.JoinPath(client.BaseURL, "/game")
	descUrl, descUrlErr := url.JoinPath(client.BaseURL, "/game/desc")
	if urlErr != nil {
		log.Println(urlErr)
	} else if descUrlErr != nil {
		log.Println(descUrlErr)
	}
	req, reqErr := http.NewRequest("GET", statusUrl, nil)
	req.Header.Set("X-Auth-Token", client.Token)
	if reqErr != nil {
		log.Println(reqErr)
	}
	descReq, descReqErr := http.NewRequest("GET", descUrl, nil)
	descReq.Header.Set("X-Auth-Token", client.Token)
	if descReqErr != nil {
		log.Println(descReqErr)
	}

	res, resErr := client.HTTPClient.Do(req)
	if resErr != nil {
		log.Println(resErr)
	}
	descRes, descResErr := client.HTTPClient.Do(descReq)
	if descResErr != nil {
		log.Println(descResErr)
	}
	defer descRes.Body.Close()
	defer res.Body.Close()

	data, dataErr := io.ReadAll(res.Body)
	if dataErr != nil {
		log.Println(dataErr)
	}
	jsonErr := json.Unmarshal(data, &stats)
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	descData, descDataErr := io.ReadAll(descRes.Body)
	if descDataErr != nil {
		log.Println(descDataErr)
	}
	jsonErr = json.Unmarshal(descData, &stats)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	return &stats, resErr
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
