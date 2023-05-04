package client

import (
	"Ships/models"
	"bytes"
	"encoding/json"
	"errors"
	gui "github.com/grupawp/warships-gui/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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
func (client *Client) Descriptions() (*models.Description, error) {
	desc := models.Description{}

	descUrl, urlErr := url.JoinPath(client.BaseURL, "/game/desc")
	if urlErr != nil {
		log.Println(urlErr)
	}
	req, reqErr := http.NewRequest("GET", descUrl, nil)
	req.Header.Set("X-Auth-Token", client.Token)
	if reqErr != nil {
		log.Println(reqErr)
	}
	time.Sleep(time.Second * 1)

	res, resErr := client.HTTPClient.Do(req)
	if resErr != nil {
		log.Println(resErr)
	}
	defer res.Body.Close()

	data, dataErr := io.ReadAll(res.Body)
	if dataErr != nil {
		log.Println(dataErr)
	}
	jsonErr := json.Unmarshal(data, &desc)
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	return &desc, resErr
}

func (client *Client) Status() (*models.Status, error) {
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
	jsonErr := json.Unmarshal(data, &stats)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	return &stats, resErr
}

func (client *Client) Fire(cord string) (string, error) {
	cord = strings.ToUpper(cord)
	isValid := isValidCoordinate(cord)
	if isValid == false {
		return "", errors.New("wrong type of coordinates. ex. (A1)")
	}
	fireData := models.Fire{Coord: cord}
	fireResp := models.FireResponse{}
	body, bodyErr := json.Marshal(fireData)
	if bodyErr != nil {
		log.Println(bodyErr)
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
	res, resErr := client.HTTPClient.Do(req)
	if resErr != nil {
		log.Println(resErr)
	}
	defer res.Body.Close()
	data, dataErr := io.ReadAll(res.Body)
	if dataErr != nil {
		log.Println(dataErr)
	}
	jsonErr := json.Unmarshal(data, &fireResp)
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	return fireResp.Result, nil
}

func isValidCoordinate(coordinate string) bool {
	pattern := `^[A-J](10|[1-9])$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(coordinate)
}
