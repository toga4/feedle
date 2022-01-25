package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

type Parser struct {
	mapper ParseItemMapper
}

func NewParser() *Parser {
	mapper := &defaultParseItemMapper{}
	return &Parser{mapper}
}

func (p *Parser) ParseURL(ctx context.Context, url string) ([]*Item, error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	return p.mapper.MapItemsFromGofeed(feed), nil
}
