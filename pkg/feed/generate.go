package feed

func (f *Feed) GenerateAtom() (string, error) {
	mapper := &defaultGenerateFeedMapper{}

	feed := mapper.MapToGorillaFeeds(f)
	atom, err := feed.ToAtom()
	if err != nil {
		return "", err
	}

	return atom, nil
}
