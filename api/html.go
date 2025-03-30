package api

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetTotalPage(body []byte) int {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	totalPage := 0
	var pageURL []string
	doc.Find("ul.pageNav-main a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			pageURL = append(pageURL, href)
		}
	})
	lastPagURL := pageURL[len(pageURL)-1]
	lastParts := lastPagURL[strings.LastIndex(lastPagURL, "/")+1:]
	fmt.Sscanf(strings.TrimPrefix(lastParts, "page-"), "%d", &totalPage)
	return totalPage + 1
}

func GetPostData(body []byte) string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	var result string
	doc.Find("div.p-body-header").Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err != nil {
			fmt.Println(err)
			return
		}
		result += html
	})
	doc.Find("div.p-body-main").Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}
		result += html
	})
	return result
}

func ScanURLPost(content *[]byte, domain string) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*content))
	if err != nil {
		return err
	}
	htmlContent := string(*content)
	cache := make(map[string]string)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			if isGotoPost(href, domain) {
				findURL := href
				if !strings.HasPrefix(findURL, fmt.Sprintf("https://%s", domain)) {
					findURL = fmt.Sprintf("https://%s/%s", domain, findURL)
				}
				fmt.Println("Found URL:", findURL)
				if cache[findURL] != "" {
					htmlContent = strings.ReplaceAll(htmlContent, href, cache[findURL])
					fmt.Println("Replace URL:", href, "to:", cache[findURL])
				} else {
					redirect := Location(findURL, findURL)
					if redirect != "" {
						htmlContent = strings.ReplaceAll(htmlContent, href, redirect)
						fmt.Println("Replace URL:", href, "to:", redirect)
						cache[findURL] = redirect
					}
				}
			}
		}
	})
	*content = []byte(htmlContent)
	return nil
}

func isGotoPost(url, domain string) bool {
	parts := strings.Split(url, "/")
	if strings.HasPrefix(url, fmt.Sprintf("https://%s/p/", domain)) || strings.HasPrefix(url, fmt.Sprintf("https://%s/goto/post?id=", domain)) || strings.HasPrefix(parts[len(parts)-1], "post-") {
		return true
	}
	return false
}
