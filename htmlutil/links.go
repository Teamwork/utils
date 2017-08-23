package htmlutil

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FindURLs will return all of the URLs in a given HTML string.
func FindURLs(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a").Each(func(idx int, s *goquery.Selection) {
		if url, exists := s.Attr("href"); exists && strings.TrimSpace(url) != "" {
			urls = append(urls, url)
		}
	})

	return urls, nil
}
