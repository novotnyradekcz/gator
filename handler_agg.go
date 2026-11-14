package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs> (e.g. 1s, 1m, 1h)")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", cmd.args[0], err)
	}
	if timeBetweenRequests <= 0 {
		return fmt.Errorf("time_between_reqs must be greater than zero")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No feeds to fetch")
			return
		}
		fmt.Printf("error getting next feed to fetch: %v\n", err)
		return
	}

	if err := s.db.MarkFeedFetched(ctx, feed.ID); err != nil {
		fmt.Printf("error marking feed %q as fetched: %v\n", feed.Name, err)
		return
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		fmt.Printf("error fetching feed %q: %v\n", feed.Name, err)
		return
	}

	fmt.Printf("Fetched %q (%d posts)\n", feed.Name, len(rssFeed.Channel.Item))
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("  * %s\n", item.Title)
	}
}
