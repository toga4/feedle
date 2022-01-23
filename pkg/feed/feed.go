package feed

import (
	"time"
)

type Feed struct {
	Title   string
	Link    string
	Updated time.Time
	ID      string
	Items   []*Item
}

type Item struct {
	Title     string
	Content   string
	Link      string
	Updated   time.Time
	Published time.Time
	Authors   []*Person
	GUID      string
}

type Person struct {
	Name  string
	Email string
}
