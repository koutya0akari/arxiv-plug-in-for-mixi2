package arxiv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const BaseURL = "https://arxiv.org"

type Paper struct {
	ID       string
	Title    string
	Authors  []string
	Category string
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func FetchNew(ctx context.Context, client HTTPClient, category string) ([]Paper, error) {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	url := fmt.Sprintf("%s/list/%s/new", BaseURL, category)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "arxiv-plug-in-for-mixi2/0.1 (+https://github.com/koutya0akari/arxiv-plug-in-for-mixi2)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: unexpected status %s", category, resp.Status)
	}

	return ParseNew(resp.Body, category)
}

func ParseNew(r io.Reader, category string) ([]Paper, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var papers []Paper
	doc.Find("h3").EachWithBreak(func(_ int, h3 *goquery.Selection) bool {
		if !strings.Contains(normalizeSpace(h3.Text()), "New submissions") {
			return true
		}
		papers = append(papers, parseEntriesAfterHeading(h3, category)...)
		return false
	})

	return papers, nil
}

func parseEntriesAfterHeading(h3 *goquery.Selection, category string) []Paper {
	var papers []Paper
	for node := h3.Next(); node.Length() > 0; node = node.Next() {
		if goquery.NodeName(node) == "h3" {
			break
		}
		switch goquery.NodeName(node) {
		case "dt":
			if paper, ok := parseDT(node, category); ok {
				papers = append(papers, paper)
			}
		case "dl":
			papers = append(papers, parseDL(node, category)...)
		default:
			node.Find("dl").Each(func(_ int, dl *goquery.Selection) {
				papers = append(papers, parseDL(dl, category)...)
			})
		}
	}
	return papers
}

func parseDL(dl *goquery.Selection, category string) []Paper {
	var papers []Paper
	dl.ChildrenFiltered("dt").Each(func(_ int, dt *goquery.Selection) {
		if paper, ok := parseDT(dt, category); ok {
			papers = append(papers, paper)
		}
	})
	return papers
}

func parseDT(dt *goquery.Selection, category string) (Paper, bool) {
	dd := dt.NextFiltered("dd")
	id := strings.TrimSpace(dt.Find(".list-identifier a[href^='/abs/']").First().Text())
	id = strings.TrimPrefix(id, "arXiv:")
	if id == "" {
		if href, ok := dt.Find("a[href^='/abs/']").First().Attr("href"); ok {
			id = strings.TrimPrefix(href, "/abs/")
		}
	}
	if id == "" || dd.Length() == 0 {
		return Paper{}, false
	}

	title := strings.TrimPrefix(normalizeSpace(dd.Find(".list-title").First().Text()), "Title:")
	title = normalizeSpace(title)
	var authors []string
	dd.Find(".list-authors a").Each(func(_ int, a *goquery.Selection) {
		author := normalizeSpace(a.Text())
		if author != "" {
			authors = append(authors, author)
		}
	})

	return Paper{
		ID:       id,
		Title:    title,
		Authors:  authors,
		Category: category,
	}, true
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
