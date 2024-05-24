package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type ExchangeRate struct {
	Bid float64 `json:"bid"`
}

func GetExchangeRate() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return 0, fmt.Errorf("error creating HTTP request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var data map[string]float64
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("error decoding JSON: %v", err)
	}

	bid := data["bid"]

	if err := SaveToFile("cotacao.txt", bid); err != nil {
		return 0, fmt.Errorf("error saving exchange rate to file: %v", err)
	}

	return bid, nil
}

func SaveToFile(filename string, bid float64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Dollar: %.2f", bid)
	if err != nil {
		return err
	}

	return nil
}
