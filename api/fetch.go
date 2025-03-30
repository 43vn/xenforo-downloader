package api

import (
	"fmt"
	"io"
	"net/http"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:136.0) Gecko/20100101 Firefox/136.0"
)

var (
	client *http.Client
)

func SetClient(c *http.Client) {
	client = c
}

func setHeaders(req *http.Request, ref string) {
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", ref)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=4")
}

func Fetch(url, ref string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	setHeaders(req, ref)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body
}

func Location(url, ref string) string {
	localClient := &http.Client{
		Jar: client.Jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", url, nil)
	setHeaders(req, ref)
	resp, err := localClient.Do(req)
	if err != nil && err != http.ErrUseLastResponse {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()
	/*
		// Debug: Print all response headers
		fmt.Println("=== Response Headers ===")
		fmt.Printf("Status: %s\n", resp.Status)
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
		fmt.Println("=======================")
	*/
	// Kiểm tra nếu có redirect (status code 3xx)
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")
		fmt.Println("Redirect location:", location)
		return location
	}

	// Nếu không có redirect
	fmt.Println("No redirect detected, status code:", resp.StatusCode)
	return ""
}
