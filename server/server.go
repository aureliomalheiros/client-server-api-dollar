package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Dollar struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func Start() error {
	fmt.Println("Initialize server on port 8080")

	db, err := sql.Open("sqlite3", "./database/cotacoes.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		getCotacao(w, r, db)
	})

	return http.ListenAndServe(":8080", nil)
}

func getCotacao(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil) // Correção aqui
	if err != nil {
		http.Error(w, "Failed to create request to external API", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get response from external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var response Dollar
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	bid := response.USDBRL.Bid

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	_, err = db.ExecContext(ctxDB, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	if err != nil {
		http.Error(w, "Failed to save to database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": bid})
}
