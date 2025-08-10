package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/mmcdole/gofeed"
)

type ParsingResult struct {
	Source      string
	Title       string
	Description string
	Link        string
}

func extractSource(feedURL string) string {
	parsed, err := url.Parse(feedURL)
	if err != nil {
		return "unknown"
	}
	host := parsed.Hostname()

	parts := strings.Split(host, ".")
	n := len(parts)
	if n < 2 {
		return host
	}
	return parts[1]
}

func main() {
	fp := gofeed.NewParser()
	//Set your rss urls
	urls := []string{}

	var wg sync.WaitGroup
	parsed := make(chan ParsingResult, 100)

	templateBytes, err := os.ReadFile("./assets/template.html")
	if err != nil {
		panic("error opening template file: " + err.Error())
	}
	templateStr := string(templateBytes)

	for _, urlStr := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			source := extractSource(u)
			feed, err := fp.ParseURL(u)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to parse %s: %v\n", u, err)
				return
			}
			for _, v := range feed.Items {
				parsed <- ParsingResult{
					Source:      source,
					Title:       v.Title,
					Description: v.Description,
					Link:        v.Link,
				}
			}
		}(urlStr)
	}

	go func() {
		wg.Wait()
		close(parsed)
	}()

	var allItems []ParsingResult
	for item := range parsed {
		allItems = append(allItems, item)
	}

	sort.Slice(allItems, func(i, j int) bool {
		if allItems[i].Source == allItems[j].Source {
			return allItems[i].Title < allItems[j].Title
		}
		return allItems[i].Source < allItems[j].Source
	})

	var itemsHTML strings.Builder
	var currentSource string
	for _, item := range allItems {
		if item.Source != currentSource {
			currentSource = item.Source
			itemsHTML.WriteString(fmt.Sprintf(
				`<h1 style="text-transform: capitalize; margin-top: 40px;">%s</h1>`,
				currentSource,
			))
		}

		itemsHTML.WriteString(`<div class="list-item">` + "\n")
		itemsHTML.WriteString(fmt.Sprintf(
			`<a href="%s" target="_blank" rel="noopener noreferrer" style="text-decoration: none; color: inherit;">`,
			item.Link,
		))
		itemsHTML.WriteString(`<h2 class="item-title">` + item.Title + "</h2>\n")
		itemsHTML.WriteString(`<p class="item-description">` + item.Description + "</p>\n")
		itemsHTML.WriteString(`</a>` + "\n")
		itemsHTML.WriteString(`</div>` + "\n")
	}

	finalHTML := strings.ReplaceAll(templateStr, "{{ITEMS}}", itemsHTML.String())
	err = os.WriteFile("result.html", []byte(finalHTML), 0644)
	if err != nil {
		panic("error writing result.html: " + err.Error())
	}
	fmt.Println("file saved")
}
