package cliente

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao")
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to get response: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	fmt.Println(string(body))

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", body))
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
		return
	}
}

