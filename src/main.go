package main

import (
	"fmt"
	"net/http"
)

func main() {
	url := "https://google.com/foo/bar"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %v\n", resp.Status)
}
