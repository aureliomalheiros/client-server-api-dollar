package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExchangeRateResponse struct {
	USDBRL ExchangeRate `json:"USDBRL"`
}

type ExchangeRate struct {
	Bid string `json:"bid"`
}

var db *sql.DB

func StartServer() error {
	var err error
	db, err = sql.Open("sqlite3", "./exchange_rates.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS exchange_rates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid REAL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	http.HandleFunc("/cotacao", HandleCotacao)

	log.Println("Servidor iniciado na porta :8080")
	return http.ListenAndServe(":8080", nil)
}

func StopServer() error {
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}

func HandleCotacao(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Println("Timeout ao chamar API de cotação")
		http.Error(w, "Timeout ao chamar API de cotação", http.StatusRequestTimeout)
		return
	default:
		resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			log.Printf("Erro ao chamar API de câmbio: %v\n", err)
			http.Error(w, "Erro ao chamar API de cotação", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Erro ao ler resposta da API: %v\n", err)
			http.Error(w, "Erro ao ler resposta da API", http.StatusInternalServerError)
			return
		}

		var data ExchangeRateResponse
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("Erro ao decodificar JSON: %v\n", err)
			http.Error(w, "Erro ao decodificar JSON", http.StatusInternalServerError)
			return
		}

		bid, err := strconv.ParseFloat(data.USDBRL.Bid, 64)
		if err != nil {
			log.Printf("Erro ao converter bid para float64: %v\n", err)
			http.Error(w, "Erro ao converter bid para float64", http.StatusInternalServerError)
			return
		}

		stmt, err := db.Prepare("INSERT INTO exchange_rates (bid) VALUES (?)")
		if err != nil {
			log.Printf("Erro ao preparar a consulta SQL: %v\n", err)
			http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(bid)
		if err != nil {
			log.Printf("Erro ao executar a consulta SQL: %v\n", err)
			http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
			return
		}

		response := map[string]float64{"bid": bid}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Erro ao codificar JSON de resposta: %v\n", err)
			http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

