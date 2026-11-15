package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/novotnyradekcz/gator/internal/database"
)

// pubDateLayouts lists the date formats commonly found in RSS/Atom feeds.
var pubDateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC3339,
	time.RFC3339Nano,
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 -0700 (MST)",
	"Mon, 2 Jan 2006 15:04 -0700",
	"Mon, 02 Jan 2006 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

func parsePublishedAt(raw string) sql.NullTime {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return sql.NullTime{}
	}
	for _, layout := range pubDateLayouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return sql.NullTime{Time: t.UTC(), Valid: true}
		}
	}
	return sql.NullTime{}
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

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

	saved := 0
	for _, item := range rssFeed.Channel.Item {
		url := strings.TrimSpace(item.Link)
		if url == "" {
			continue
		}
		title := strings.TrimSpace(item.Title)
		if title == "" {
			title = url
		}

		now := time.Now()
		_, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Title:     title,
			Url:       url,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: parsePublishedAt(item.PubDate),
			FeedID:      feed.ID,
		})
		if err != nil {
			if isUniqueViolation(err) {
				continue
			}
			fmt.Printf("error saving post %q: %v\n", url, err)
			continue
		}
		saved++
	}

	fmt.Printf("Fetched %q: %d items, %d new posts\n", feed.Name, len(rssFeed.Channel.Item), saved)
}
