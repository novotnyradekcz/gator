package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/novotnyradekcz/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("exactly a URL is required\nusage: gator follow <url>")
	}

	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("get feed by URL: %w", err)
	}

	now := time.Now()
	follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("create feed follow: %w", err)
	}

	fmt.Printf("%s is now following %s\n", follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following does not accept arguments\nusage: gator following")
	}

	ctx := context.Background()
	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("get feed follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("Not following any feeds.")
		return nil
	}

	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, f := range follows {
		fmt.Printf("* %s\n", f.FeedName)
	}
	return nil
}
