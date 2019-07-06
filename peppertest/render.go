package peppertest

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/elliotchance/pepper"
	"strings"
)

func RenderToDocument(c pepper.Component) (*goquery.Document, error) {
	html, err := pepper.Render(c)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	return doc, nil
}
