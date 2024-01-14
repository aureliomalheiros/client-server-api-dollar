package main

import (
	"fmt"
	"net/http"
	"encoding/json"
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
func main(){
	fmt.Println("Initialize my server in port 8080")
	http.HandleFunc("/cotacao", cotacao())
	http.ListenAndServe(":8080", nil)
}

func cotacao(w http.ResponseWriter, r *http.Request){
	res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	var response Dollar

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Bid: %s\n", response.USDBRL.Bid)
}