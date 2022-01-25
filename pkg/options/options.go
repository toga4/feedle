package options

import (
	"flag"
)

type Flags struct {
	Title string
	Link  string
}

func Parse(args []string) (*Flags, error) {
	flags := &Flags{}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&flags.Title, "title", "Feed bundled by feedle", "Title of bundled feed")
	fs.StringVar(&flags.Link, "link", "", "Link of bundled feed")

	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}

	return flags, nil
}
