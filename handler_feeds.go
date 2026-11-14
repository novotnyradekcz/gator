package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feeds does not accept arguments\nusage: gator feeds")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("get feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Feed: %s\n", feed.FeedName)
		fmt.Printf("URL: %s\n", feed.FeedUrl)
		fmt.Printf("User: %s\n\n", feed.UserName)
	}
	return nil
}
