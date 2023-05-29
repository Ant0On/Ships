package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
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

func NewClient(baseURL string, timeOut time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeOut,
		},
	}
}
func (client *Client) InitGame(initialProps InitialData) error {
	time.Sleep(time.Second * 1)
	initialData := initialProps

	jsonData, jsonDataErr := json.Marshal(initialData)
	if jsonDataErr != nil {
		log.Println(jsonDataErr)
		return jsonDataErr
	}
	reader := bytes.NewReader(jsonData)

	initUrl, urlErr := url.JoinPath(client.BaseURL, "/game")
	if urlErr != nil {
		log.Println(urlErr)
		return urlErr
	}

	req, reqErr := http.NewRequest("POST", initUrl, reader)
	if reqErr != nil {
		log.Println(reqErr)
		return reqErr
	}

	resp, respErr := client.HTTPClient.Do(req)
	client.handleBadReq(resp, respErr, req)

	slog.Info("client [InitGame]", slog.Any("initialData", initialData))

	if respErr != nil {
		log.Println(respErr)
		return respErr
	}
	slog.Info("client [InitGame]", slog.Int("statusCode", resp.StatusCode))
	client.Token = resp.Header.Get("X-Auth-Token")
	slog.Info("client [InitGame]", slog.String("token", client.Token))

	return nil
}
func (client *Client) Board() ([]string, error) {
	board := Board{}
	boardURL, urlErr := url.JoinPath(client.BaseURL, "/game/board")

	if urlErr != nil {
		log.Println(urlErr)
		return nil, urlErr
	}
	cords, _ := client.get(boardURL)

	jsonErr := json.Unmarshal(cords, &board)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return board.Board, nil

}

func (client *Client) Descriptions() (*Description, error) {
	data := Description{}

	descUrl, urlErr := url.JoinPath(client.BaseURL, "/game/desc")
	if urlErr != nil {
		log.Println(urlErr)
		return nil, urlErr
	}
	bodyData, _ := client.get(descUrl)
	jsonErr := json.Unmarshal(bodyData, &data)
	if jsonErr != nil {
		log.Println(jsonErr)
		return nil, jsonErr
	}

	return &data, nil
}

func (client *Client) Status() (*Status, error) {
	var status Status
	statusUrl, urlErr := url.JoinPath(client.BaseURL, "/game")
	if urlErr != nil {
		log.Println(urlErr)
		return nil, urlErr
	}
	data, _ := client.get(statusUrl)
	jsonErr := json.Unmarshal(data, &status)
	if jsonErr != nil {
		log.Println(jsonErr)
		return nil, jsonErr
	}

	return &status, nil
}

func (client *Client) Fire(cord string) (string, error) {
	isValid := isValidCoordinate(cord)
	if isValid == false {
		return "", errors.New("wrong type of coordinates. ex. (A1)")
	}
	fireData := Fire{Coord: cord}
	fireResp := FireResponse{}
	body, bodyErr := json.Marshal(fireData)
	if bodyErr != nil {
		log.Println(bodyErr)
		return "", bodyErr
	}

	fireUrl, fireErr := url.JoinPath(client.BaseURL, "/game/fire")
	if fireErr != nil {
		log.Println(fireErr)
		return "", fireErr
	}

	data, _ := client.post(fireUrl, bytes.NewReader(body))

	jsonErr := json.Unmarshal(data, &fireResp)
	if jsonErr != nil {
		log.Println(jsonErr)
		return "", jsonErr
	}
	return fireResp.Result, nil
}

func (client *Client) PlayersList() error {
	var listData []ListData
	time.Sleep(time.Second * 1)
	listUrl, urlErr := url.JoinPath(client.BaseURL, "/game/list")
	if urlErr != nil {
		log.Println(urlErr)
		return urlErr
	}
	data, _ := client.get(listUrl)
	jsonErr := json.Unmarshal(data, &listData)
	if jsonErr != nil {
		log.Println(jsonErr)
		return jsonErr
	}
	for _, d := range listData {
		fmt.Println("Nick: ", d.Nick)
	}
	return nil

}

func (client *Client) Abandon() error {
	abandonUrl, urlErr := url.JoinPath(client.BaseURL, "/game/abandon")
	if urlErr != nil {
		log.Println(urlErr)
		return urlErr
	}
	req, reqErr := http.NewRequest("DELETE", abandonUrl, nil)
	req.Header.Set("X-Auth-Token", client.Token)
	if reqErr != nil {
		log.Println(reqErr)
		return reqErr
	}
	res, resErr := client.HTTPClient.Do(req)
	client.handleBadReq(res, resErr, req)
	if resErr != nil {
		log.Println(resErr)
		return resErr
	}
	defer res.Body.Close()

	// print the response body
	return nil
}

func (client *Client) Refresh() error {
	refreshUrl, urlErr := url.JoinPath(client.BaseURL, "/game/refresh")
	if urlErr != nil {
		log.Println(urlErr)
		return urlErr
	}
	_, getErr := client.get(refreshUrl)
	if getErr != nil {
		return getErr
	}
	return nil
}

func (client *Client) Top10() error {
	var top10 Top10
	topUrl, urlErr := url.JoinPath(client.BaseURL, "/stats")
	if urlErr != nil {
		return urlErr
	}
	data, dataErr := client.get(topUrl)
	if dataErr != nil {
		return dataErr
	}
	jsonErr := json.Unmarshal(data, &top10)
	if jsonErr != nil {
		log.Println(jsonErr)
		return jsonErr
	}
	top10Data := make([][]string, 11)
	top10Data[0] = append(top10Data[0], "Rank", "Nick", "Games", "Wins", "Points")
	for i, d := range top10.Stats {
		top10Data[i+1] = append(top10Data[i+1], fmt.Sprintf("%d", d.Rank), d.Nick,
			fmt.Sprintf("%d", d.Games), fmt.Sprintf("%d", d.Wins), fmt.Sprintf("%d", d.Points))
	}
	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(top10Data).Render()

	return nil
}

func (client *Client) PlayerStats(player string) error {
	var playerStats PlayerStats
	topUrl, urlErr := url.JoinPath(client.BaseURL, fmt.Sprintf("/stats/%s", player))
	if urlErr != nil {
		return urlErr
	}
	data, dataErr := client.get(topUrl)
	if dataErr != nil {
		return dataErr
	}
	jsonErr := json.Unmarshal(data, &playerStats)
	if jsonErr != nil {
		log.Println(jsonErr)
		return jsonErr
	}

	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
		{"Rank", "Nick", "Games", "Wins", "Points"},
		{fmt.Sprintf("%d", playerStats.Stats.Rank), playerStats.Stats.Nick,
			fmt.Sprintf("%d", playerStats.Stats.Games), fmt.Sprintf("%d", playerStats.Stats.Wins),
			fmt.Sprintf("%d", playerStats.Stats.Points)},
	}).Render()

	return nil
}

func isValidCoordinate(coordinate string) bool {
	pattern := `^[A-J](10|[1-9])$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(coordinate)
}
func (client *Client) get(url string) ([]byte, error) {
	req, reqErr := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", client.Token)
	if reqErr != nil {
		log.Println(reqErr)
		return nil, reqErr
	}
	time.Sleep(time.Millisecond * 100)

	res, resErr := client.HTTPClient.Do(req)

	if resErr != nil {
		log.Println(resErr)
		return nil, resErr
	}
	client.handleBadReq(res, resErr, req)
	defer res.Body.Close()

	bodyData, dataErr := io.ReadAll(res.Body)
	if dataErr != nil {
		log.Println(dataErr)
		return nil, dataErr
	}
	return bodyData, nil
}

func (client *Client) post(url string, reader *bytes.Reader) ([]byte, error) {
	req, reqErr := http.NewRequest("POST", url, reader)
	req.Header.Set("X-Auth-Token", client.Token)
	if reqErr != nil {
		log.Println(reqErr)
		return nil, reqErr
	}
	time.Sleep(time.Millisecond * 100)
	res, resErr := client.HTTPClient.Do(req)
	time.Sleep(time.Millisecond * 200)
	client.handleBadReq(res, resErr, req)
	if resErr != nil {
		log.Println(resErr)
		return nil, resErr
	}
	defer res.Body.Close()
	data, dataErr := io.ReadAll(res.Body)
	if dataErr != nil {
		log.Println(dataErr)
		return nil, dataErr
	}
	return data, nil
}
func (client *Client) handleBadReq(resp *http.Response, respErr error, req *http.Request) {
	for resp.StatusCode == 503 {
		time.Sleep(time.Second * 1)
		resp, respErr = client.HTTPClient.Do(req)
	}
}
