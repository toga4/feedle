package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"
	"github.com/toga4/github-stars-release-feeder/pkg/iterator"
)

type Client struct {
	client *github.Client
}

type Repository struct {
	*github.Repository
}

func NewClient() *Client {
	c := github.NewClient(nil)
	return &Client{c}
}

func (c *Client) ListStarred(ctx context.Context, user string) *RepositoryIterator {
	page := 1
	return &RepositoryIterator{
		fetcher: func() ([]*Repository, bool, error) {
			repos, resp, err := c.client.Activity.ListStarred(ctx, user, &github.ActivityListStarredOptions{
				ListOptions: github.ListOptions{
					Page: page,
				},
			})
			if err != nil {
				return nil, false, fmt.Errorf("github.Repository#ListStarred error: %w", err)
			}

			repositories := []*Repository{}
			for _, repo := range repos {
				r := &Repository{repo.Repository}
				repositories = append(repositories, r)
			}

			page = resp.NextPage
			last := resp.NextPage == resp.LastPage

			return repositories, last, nil
		},
	}
}

func MakeReleaseFeedUrl(fullName string) string {
	return fmt.Sprintf("https://github.com/%s/releases.atom", fullName)
}

type RepositoryIterator struct {
	fetcher func() ([]*Repository, bool, error)

	err   error
	last  bool
	repos []*Repository
}

func (ri *RepositoryIterator) Next() (*Repository, error) {
	if ri.err != nil {
		return nil, ri.err
	}

	if len(ri.repos) == 0 && !ri.last {
		ri.repos, ri.last, ri.err = ri.fetcher()
		if ri.err != nil {
			return nil, ri.err
		}
	}

	if len(ri.repos) > 0 {
		row := ri.repos[0]
		ri.repos = ri.repos[1:]
		return row, nil
	}

	return nil, iterator.ErrDone
}
