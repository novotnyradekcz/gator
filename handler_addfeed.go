package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/novotnyradekcz/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("exactly a name and URL are required\nusage: gator addfeed <name> <url>")
	}

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("get current user: %w", err)
	}

	now := time.Now()
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("create feed: %w", err)
	}

	fmt.Println("Feed created:")
	fmt.Printf("  ID: %s\n", feed.ID)
	fmt.Printf("  Created At: %s\n", feed.CreatedAt)
	fmt.Printf("  Updated At: %s\n", feed.UpdatedAt)
	fmt.Printf("  Name: %s\n", feed.Name)
	fmt.Printf("  URL: %s\n", feed.Url)
	fmt.Printf("  User ID: %s\n", feed.UserID)
	return nil
}
