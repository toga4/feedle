package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/toga4/gsrf/pkg/feed"
	"github.com/toga4/gsrf/pkg/github"
	"github.com/toga4/gsrf/pkg/iterator"
)

func main() {
	ctx := context.Background()

	flags := struct {
		User string
		Link string
	}{}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&flags.User, "user", "", "Github username")
	fs.StringVar(&flags.Link, "link", "", "URL of feed")
	fs.Parse(os.Args[1:])

	githubClient := github.NewClient()
	parser := feed.NewParser()

	feed := &feed.Feed{}
	feed.Title = fmt.Sprintf("Releases from Starred by %s", flags.User)
	feed.Link = flags.Link
	feed.ID = flags.Link

	iter := githubClient.ListStarred(ctx, flags.User)

	repoCh := NextRepositries(ctx, iter)
	feed.Items = ParseFeedItems(ctx, parser, repoCh)

	log.Println("Sorting...")

	sort.Slice(feed.Items, func(i, j int) bool {
		return feed.Items[i].Updated.After(feed.Items[j].Updated)
	})

	log.Println("Generating feed...")
	feed.Updated = feed.Items[0].Updated
	atom, err := feed.GenerateAtom()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(atom)
}

func NextRepositries(ctx context.Context, iter *github.RepositoryIterator) <-chan *github.Repository {
	ch := make(chan *github.Repository)

	go func() {
		defer close(ch)

	LOOP:
		for {
			repo, err := iter.Next()
			if err != nil {
				if err == iterator.ErrDone {
					break
				}
				log.Fatal(err)
			}

			select {
			case <-ctx.Done():
				break LOOP
			case ch <- repo:
			}
		}

	}()

	return ch
}

func ParseFeedItems(ctx context.Context, parser *feed.Parser, ch <-chan *github.Repository) []*feed.Item {
	feedItems := []*feed.Item{}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for repo := range ch {

		repo := repo

		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Printf("Fetch release feed %s\n", *repo.FullName)
			releaseFeedUrl := github.MakeReleaseFeedUrl(*repo.FullName)

			items, err := parser.ParseURL(ctx, releaseFeedUrl)
			if err != nil {
				log.Fatal(err)
			}

			mu.Lock()
			defer mu.Unlock()
			feedItems = append(feedItems, items...)
		}()
	}

	wg.Wait()

	return feedItems
}
