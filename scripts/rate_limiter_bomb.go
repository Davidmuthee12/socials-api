package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	url := os.Getenv("RATE_LIMITER_BOMB_URL")
	if url == "" {
		url = "http://127.0.0.1:8080/v1/users/feed"
	}

	requests := 25
	if value := os.Getenv("RATE_LIMITER_BOMB_REQUESTS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			fmt.Printf("invalid RATE_LIMITER_BOMB_REQUESTS %q\n", value)
			os.Exit(1)
		}
		requests = parsed
	}

	for i := 0; i < requests; i++ {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Request %d -> %d\n", i+1, resp.StatusCode)
		resp.Body.Close()
	}
}
