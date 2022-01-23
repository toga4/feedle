package feed

import (
	"context"

	"github.com/mmcdole/gofeed"
)

type Parser struct {
	parser *gofeed.Parser
	mapper ParseItemMapper
}

func NewParser() *Parser {
	parser := gofeed.NewParser()
	mapper := &defaultParseItemMapper{}
	return &Parser{parser, mapper}
}

func (p *Parser) ParseURL(ctx context.Context, url string) ([]*Item, error) {
	feed, err := p.parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	return p.mapper.MapItemsFromGofeed(feed), nil
}
