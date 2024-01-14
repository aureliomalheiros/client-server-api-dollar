package main

import (
	"fmt"
	"net/http"
	"context"
	"time"
)
func main(){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Milliseconds)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080")

	if err != nil {
		panic(err)
	}

}