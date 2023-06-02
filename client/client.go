package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pterm/pterm"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"net/http"
	"net/url"
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
	initialData := initialProps

	jsonData, err := json.Marshal(initialData)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(jsonData)

	initUrl, err := url.JoinPath(client.BaseURL, "/game")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", initUrl, reader)
	if err != nil {
		return err
	}

	clientCopy := client.createClient()

	resp, err := clientCopy.Do(req)

	slog.Info("client [InitGame]", slog.Any("initialData", initialData))

	if err != nil {
		return err
	}
	slog.Info("client [InitGame]", slog.Int("statusCode", resp.StatusCode))
	client.Token = resp.Header.Get("X-Auth-Token")
	slog.Info("client [InitGame]", slog.String("token", client.Token))

	return nil
}
func (client *Client) Board() ([]string, error) {
	board := Board{}
	boardURL, err := url.JoinPath(client.BaseURL, "/game/board")

	if err != nil {
		return nil, err
	}
	cords, _ := client.get(boardURL)

	err = json.Unmarshal(cords, &board)
	if err != nil {
		return nil, err
	}
	return board.Board, nil

}

func (client *Client) Descriptions() (*Description, error) {
	data := Description{}

	descUrl, err := url.JoinPath(client.BaseURL, "/game/desc")
	if err != nil {
		return nil, err
	}
	bodyData, _ := client.get(descUrl)
	err = json.Unmarshal(bodyData, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (client *Client) Status() (*Status, error) {
	var status Status
	statusUrl, err := url.JoinPath(client.BaseURL, "/game")
	if err != nil {
		return nil, err
	}
	data, _ := client.get(statusUrl)
	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (client *Client) Fire(cord string) (string, error) {

	fireData := Fire{Coord: cord}
	fireResp := FireResponse{}
	body, err := json.Marshal(fireData)
	if err != nil {
		log.Fatal(err)
	}

	fireUrl, err := url.JoinPath(client.BaseURL, "/game/fire")
	if err != nil {
		log.Fatal(err)
	}

	data, _ := client.post(fireUrl, bytes.NewReader(body))

	err = json.Unmarshal(data, &fireResp)
	if err != nil {
		log.Fatal(err)
	}
	return fireResp.Result, nil
}

func (client *Client) PlayersList() error {
	var listData []ListData
	listUrl, err := url.JoinPath(client.BaseURL, "/game/list")
	if err != nil {
		return err
	}
	data, _ := client.get(listUrl)
	err = json.Unmarshal(data, &listData)
	if err != nil {
		return err
	}

	return nil

}

func (client *Client) Abandon() error {
	abandonUrl, err := url.JoinPath(client.BaseURL, "/game/abandon")
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", abandonUrl, nil)
	req.Header.Set("X-Auth-Token", client.Token)
	if err != nil {
		return err
	}

	clientCopy := client.createClient()
	res, err := clientCopy.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (client *Client) Refresh() error {
	refreshUrl, err := url.JoinPath(client.BaseURL, "/game/refresh")
	if err != nil {
		return err
	}
	_, err = client.get(refreshUrl)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Top10() error {
	var top10 Top10
	topUrl, err := url.JoinPath(client.BaseURL, "/stats")
	if err != nil {
		return err
	}
	data, err := client.get(topUrl)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &top10)
	if err != nil {
		return err
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
	topUrl, err := url.JoinPath(client.BaseURL, fmt.Sprintf("/stats/%s", player))
	if err != nil {
		return err
	}
	data, err := client.get(topUrl)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &playerStats)
	if err != nil {
		return err
	}

	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(pterm.TableData{
		{"Rank", "Nick", "Games", "Wins", "Points"},
		{fmt.Sprintf("%d", playerStats.Stats.Rank), playerStats.Stats.Nick,
			fmt.Sprintf("%d", playerStats.Stats.Games), fmt.Sprintf("%d", playerStats.Stats.Wins),
			fmt.Sprintf("%d", playerStats.Stats.Points)},
	}).Render()

	return nil
}

func (client *Client) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", client.Token)
	if err != nil {
		return nil, err
	}

	clientCopy := client.createClient()

	res, err := clientCopy.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bodyData, nil
}

func (client *Client) post(url string, reader *bytes.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", url, reader)
	req.Header.Set("X-Auth-Token", client.Token)
	if err != nil {
		return nil, err
	}

	clientCopy := client.createClient()

	res, err := clientCopy.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (client *Client) createClient() *http.Client {
	client2 := retryablehttp.NewClient()
	client2.HTTPClient = client.HTTPClient
	client2.RetryMax = 10
	client2.RetryWaitMin = time.Second
	client2.Logger = nil
	client3 := client2.StandardClient()
	return client3
}
