package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/toga4/feedle/pkg/feed"
	"github.com/toga4/feedle/pkg/options"
)

func main() {
	ctx := context.Background()

	flags, err := options.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	parser := feed.NewParser()

	repoCh := GenerateChannel(ctx, os.Stdin)

	itemsChannels := make([]<-chan []*feed.Item, 0)
	for i := 0; i < 10; i++ {
		itemsCh := ParseFeedItems(ctx, parser, repoCh)
		itemsChannels = append(itemsChannels, itemsCh)
	}

	items := []*feed.Item{}
	for is := range merge(itemsChannels...) {
		items = append(items, is...)
	}

	log.Println("Generating feed...")

	sort.Slice(items, func(i, j int) bool {
		return items[i].Updated.After(items[j].Updated)
	})

	feed := &feed.Feed{}
	feed.Title = flags.Title
	feed.Link = flags.Link
	feed.ID = flags.Link
	feed.Updated = items[0].Updated
	feed.Items = items

	atom, err := feed.GenerateAtom()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(atom)
}

func GenerateChannel(ctx context.Context, reader io.Reader) <-chan string {
	out := make(chan string)

	in := bufio.NewScanner(reader)

	go func() {
		defer close(out)

	LOOP:
		for in.Scan() {
			select {
			case <-ctx.Done():
				break LOOP
			case out <- in.Text():
			}
		}

		// handle first encountered error while reading
		if err := in.Err(); err != nil {
			log.Fatalf("Error while reading input: %v", err)
		}
	}()

	return out
}

func ParseFeedItems(ctx context.Context, parser *feed.Parser, urlCh <-chan string) <-chan []*feed.Item {
	out := make(chan []*feed.Item)

	go func() {
		defer close(out)
		for url := range urlCh {
			log.Printf("Fetch: %s\n", url)
			items, err := parser.ParseURL(ctx, url)
			if err != nil {
				log.Fatalf("Error while parsing feed: %s, err: %v", url, err)
			}
			log.Printf("Fetched: %s\n", url)

			out <- items
		}
	}()

	return out
}

func merge(cs ...<-chan []*feed.Item) <-chan []*feed.Item {
	out := make(chan []*feed.Item)

	var wg sync.WaitGroup

	for _, c := range cs {
		wg.Add(1)
		go func(c <-chan []*feed.Item) {
			defer wg.Done()
			for n := range c {
				out <- n
			}
		}(c)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
