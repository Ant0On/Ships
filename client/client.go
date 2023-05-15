package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type IClient interface {
	InitGame(desc, nick string, wpBot bool) (*InitialData, error)
	Board() ([]string, error)
	Descriptions() (*Description, error)
	Status() (*Status, error)
	Fire(cord string) (string, error)
}

func NewClient(baseURL string, timeOut time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeOut,
		},
	}
}
func (client *Client) InitGame(nick, desc, targetNick string, wpBot bool) error {
	time.Sleep(time.Second * 1)
	initialData := InitialData{
		Desc:       desc,
		Nick:       nick,
		TargetNick: targetNick,
		Wpbot:      wpBot,
	}

	jsonData, jsonDataErr := json.Marshal(initialData)
	if jsonDataErr != nil {
		log.Println(jsonDataErr)
		return jsonDataErr
	}
	reader := bytes.NewReader(jsonData)

	initUrl, urlErr := url.JoinPath(client.BaseURL, "/game")
	if urlErr != nil {
		log.Println(urlErr)
	}

	req, reqErr := http.NewRequest(http.MethodPost, initUrl, reader)
	if reqErr != nil {
		log.Println(reqErr)
		return reqErr
	}

	resp, respErr := client.HTTPClient.Do(req)

	slog.Info("client [InitGame]", slog.Any("initialData", initialData))

	if respErr != nil {
		log.Println(respErr)
		return respErr
	}
	slog.Info("client [InitGame]", slog.Int("statusCode", resp.StatusCode))
	client.Token = resp.Header.Get("X-Auth-Token")
	slog.Info("client [InitGame]", slog.String("token", client.Token))

	return respErr
}
func (client *Client) Board() ([]string, error) {
	board := Board{}
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

func (client *Client) Descriptions() (*Description, error) {
	desc := Description{}

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

func (client *Client) Status() (*Status, error) {
	var status Status
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
	jsonErr := json.Unmarshal(data, &status)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	return &status, resErr
}

func (client *Client) Fire(cord string) (string, error) {
	isValid := IsValidCoordinate(cord)
	if isValid == false {
		return "", errors.New("wrong type of coordinates. ex. (A1)")
	}
	fireData := Fire{Coord: cord}
	fireResp := FireResponse{}
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

func (client *Client) PlayersList() (ListData, error) {
	var listData ListData
	time.Sleep(time.Second * 1)
	listUrl, urlErr := url.JoinPath(client.BaseURL, "/game/list")
	if urlErr != nil {
		log.Println(urlErr)
		return listData, urlErr
	}
	req, reqErr := http.NewRequest("GET", listUrl, nil)
	if reqErr != nil {
		log.Println(reqErr)
		return listData, reqErr
	}
	res, resErr := client.HTTPClient.Do(req)
	if resErr != nil {
		log.Println(resErr)
		return listData, resErr
	}
	defer res.Body.Close()
	data, dataErr := io.ReadAll(res.Body)
	if dataErr != nil {
		log.Println(dataErr)
	}
	jsonErr := json.Unmarshal(data, &listData)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	return listData, resErr
}

func IsValidCoordinate(coordinate string) bool {
	pattern := `^[A-J](10|[1-9])$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(coordinate)
}
