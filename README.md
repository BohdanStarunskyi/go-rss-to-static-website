# Go RSS To Static Website

A simple, concurrent RSS feed parser written in Go. Fetches articles from multiple news sources, groups them by source, and generates a styled HTML file you can open in your browser.

## Project Description

* Parses multiple RSS feeds concurrently for speed.
* Automatically extracts the source name from the feed URL.
* Groups articles by source and sorts them alphabetically.
* Generates a clickable HTML list of articles using the included HTML template (`assets/template.html`).
* The template is fully customizable — you can change its layout, colors, and styling to match your needs.

## Features

* Concurrent feed fetching with goroutines and a `sync.WaitGroup`.
* Parsing powered by the `gofeed` library.
* Source extraction using the feed URL's hostname.
* Output HTML (`result.html`) created by replacing a `{{ITEMS}}` placeholder in a template file.

## Video demo
[Link](https://youtu.be/f7fh2Nartok)

## Tech Stack

* **Go** (Golang)
* **gofeed** — [https://github.com/mmcdole/gofeed](https://github.com/mmcdole/gofeed) for RSS parsing
* Standard Go libraries for concurrency, I/O, and string manipulation

## Project Structure

```
assets/             # Template files (e.g., template.html)
main.go             # Entry point: RSS fetching, parsing, and HTML generation
result.html         # Generated HTML file with all feed items
go.mod / go.sum     # Dependency management
```

## How it Works (high level)

1. The program loads `assets/template.html` and finds the `{{ITEMS}}` placeholder.
2. For each URL in the `urls` slice it spawns a goroutine that:

   * Parses the RSS feed using `gofeed`.
   * Extracts article `Title`, `Description`, and `Link`.
   * Sends results back over a buffered channel.
3. A collector goroutine waits for all feed goroutines to finish, closes the channel, and aggregates items into `allItems`.
4. Items are sorted and grouped by `Source` and rendered into the template.
5. The final `result.html` file is written to disk.

## Example `urls` slice

```go
urls := []string{
    "https://feeds.bbci.co.uk/news/world/rss.xml",
    "http://rss.cnn.com/rss/edition.rss",
    "https://www.wired.com/feed/rss",
}
```

## Template

The project already includes a `template.html` in the `assets/` folder. It contains a `{{ITEMS}}` placeholder where the generated content will be inserted. You can freely modify the template's HTML or CSS to adjust the page design.

Minimal example:

```html
<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>RSS Aggregator</title>
    <style>
      body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial; margin:24px; background:#f7f8fb; color:#111 }
      .list-item { padding:12px; border-bottom:1px solid #eee }
      h1 { text-transform: capitalize; margin-top:36px }
      .item-title { margin:0 0 6px 0 }
      .item-description { margin:0; color:#444 }
      a { color: inherit; text-decoration:none }
    </style>
  </head>
  <body>
    <h1>Aggregated Feeds</h1>
    {{ITEMS}}
  </body>
</html>
```

## How to Run

1. **Install Go** — [https://golang.org/dl/](https://golang.org/dl/)
2. **Clone the repository**

```bash
git clone https://github.com/your-username/go-rss-to-static-website.git
cd go-rss-to-static-website
```

3. **Add your RSS feed URLs** — edit the `urls` slice in `main.go`.
4. **Download dependencies**

```bash
go mod tidy
```

5. **Run**

```bash
go run main.go
```

6. **Open** `result.html` in your browser.

## Notes & Tips

* The `extractSource` function uses the hostname parts to derive a short source name; adjust it if you need full domain names.
* Increase the channel buffer if you parse many feeds simultaneously.
* Consider sanitizing or truncating long descriptions when rendering HTML.
* For production use, add error handling/logging and consider rate-limiting or caching.

## License

MIT — feel free to reuse and modify.
