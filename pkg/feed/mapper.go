package feed

import (
	"fmt"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

type ParseItemMapper interface {
	MapItemsFromGofeed(f *gofeed.Feed) []*Item
}

type defaultParseItemMapper struct{}

func (*defaultParseItemMapper) MapItemsFromGofeed(feed *gofeed.Feed) []*Item {
	feedTitle := feed.Title

	items := []*Item{}
	for _, item := range feed.Items {
		title := fmt.Sprintf("%s - %s", feedTitle, item.Title)

		authors := []*Person{}
		for _, author := range item.Authors {
			authors = append(authors, &Person{
				Name:  author.Name,
				Email: author.Email,
			})
		}

		items = append(items, &Item{
			Title:     title,
			Content:   item.Content,
			Link:      item.Link,
			Updated:   *item.UpdatedParsed,
			Published: *item.PublishedParsed,
			Authors:   authors,
			GUID:      item.GUID,
		})
	}

	return items
}

type GenerateFeedMapper interface {
	MapToGorillaFeeds(f *Feed) *feeds.Feed
}

type defaultGenerateFeedMapper struct{}

func (*defaultGenerateFeedMapper) MapToGorillaFeeds(feed *Feed) *feeds.Feed {
	items := []*feeds.Item{}
	for _, item := range feed.Items {
		items = append(items, &feeds.Item{
			Title:   item.Title,
			Link:    &feeds.Link{Href: item.Link},
			Content: item.Content,
			Updated: item.Updated,
			Created: item.Published,
			Author:  &feeds.Author{Name: item.Authors[0].Name, Email: item.Authors[0].Email},
			Id:      item.GUID,
		})
	}

	return &feeds.Feed{
		Title:   feed.Title,
		Link:    &feeds.Link{Href: feed.Link},
		Updated: feed.Updated,
		Id:      feed.ID,
		Items:   items,
	}
}
