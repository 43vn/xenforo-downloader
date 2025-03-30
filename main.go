package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
	"xenforo-downloader/api"

	"golang.org/x/net/proxy"
)

var (
	Url, ref, toc, scan, domain, Proxy, output string
	parts                                      []string
)

func init() {
	transport := &http.Transport{}
	if Proxy != "" {
		if strings.HasPrefix(Proxy, "socks5://") {
			dialer, err := proxy.SOCKS5("tcp", strings.TrimPrefix(Proxy, "socks5://"), nil, proxy.Direct) // Replace with your SOCKS5 proxy address
			if err != nil {
				log.Fatal("Error creating SOCKS5 dialer:", err)
			}
			transport = &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
			}
		} else if strings.HasPrefix(Proxy, "http://") {
			proxyURL, err := url.Parse(Proxy)
			if err != nil {
				log.Fatal("Error parsing proxy URL:", err)
			}
			transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		}
	}
	flag.StringVar(&Url, "url", "", "URL of the website")
	flag.StringVar(&ref, "ref", "", "Referer URL")
	flag.StringVar(&toc, "toc", "", "File JSON Table of Content")
	flag.StringVar(&scan, "scan", "", "Scan Path")
	flag.StringVar(&domain, "domain", "voz.vn", "Domain")
	flag.StringVar(&Proxy, "proxy", "", "Proxy")
	flag.StringVar(&output, "output", "", "Output directory")
	flag.Parse()
	if Url == "" && toc == "" && scan == "" {
		log.Fatal("Please provide a URL or TOC file or Scan Path")
	} else {
		parts = strings.Split(Url, "/")
	}
	if ref == "" {
		ref = Url
	}
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal("Error creating cookie jar:", err)
	}

	api.SetClient(&http.Client{
		Jar:       cookieJar,
		Transport: transport,
	})
}

func main() {
	if Url != "" {
		processURL()
	} else if toc != "" {
		processTOC()
	} else if scan != "" {
		processScan()
	}
}

func processTOC() {
	content, err := os.ReadFile(toc)
	if err != nil {
		log.Fatal("Error reading TOC file:", err)
	}
	var links []string
	err = json.Unmarshal(content, &links)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}
	result := map[string]string{}
	for _, link := range links {
		if ref == "" {
			ref = link
		}
		fmt.Println("Checking URL: ", link, "with referer:", ref)
		redirect := api.Location(link, ref)
		if redirect != "" {
			result[link] = redirect
		}
		time.Sleep(2 * time.Second)
	}
	jsonData, _ := json.Marshal(result)
	os.WriteFile(strings.ReplaceAll(toc, "toc.json", "direct.json"), jsonData, 0644)
	fmt.Printf("Done with %d links", len(result))
}
func processURL() {
	body := api.Fetch(Url, ref)
	totalPage := api.GetTotalPage(body)
	indexContent := api.GetPostData(body)
	if totalPage > 0 && indexContent != "" {
		projectName := parts[len(parts)-2]
		if output != "" {
			projectName = output
		}
		createDir(projectName)
		if !isExists(fmt.Sprintf("%s/index.html", projectName)) {
			os.WriteFile(fmt.Sprintf("%s/index.html", projectName), []byte(indexContent), 0644)
		}
		for p := range totalPage {
			if p == 0 || p == 1 {
				continue
			}
			fileName := fmt.Sprintf("%s/page-%d.html", projectName, p)
			if isExists(fileName) {
				continue
			}
			for retries := range 16 {
				body := api.Fetch(fmt.Sprintf("%spage-%d", Url, p), ref)
				if body != nil {
					pageContent := api.GetPostData(body)
					if pageContent != "" {
						os.WriteFile(fileName, []byte(pageContent), 0644)
						break
					}
				}
				if retries == 15 {
					fmt.Printf("Failed to fetch page %d after 16 retries\n", p)
				}
				time.Sleep(time.Second * time.Duration(retries) * 2)
			}
			time.Sleep(time.Second * 1)
		}
	} else {
		fmt.Println("No data found")
	}
}

func processScan() {
	d, err := os.ReadDir(scan)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range d {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".html") {
			fmt.Println("Checking file:", e.Name())
			content, err := os.ReadFile(fmt.Sprintf("%s/%s", scan, e.Name()))
			if err != nil {
				fmt.Println("File:", e.Name(), "Lỗi:", err)
				continue
			}
			err = api.ScanURLPost(&content, domain)
			if err == nil {
				os.WriteFile(fmt.Sprintf("%s/%s", scan, e.Name()), []byte(content), 0644)
			} else {
				fmt.Println("File:", e.Name(), "Lỗi:", err)
			}
		}
	}
}

func createDir(projectName string) {
	if _, err := os.Stat(projectName); os.IsNotExist(err) {
		os.Mkdir(projectName, 0755)
	}
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
