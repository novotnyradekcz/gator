package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/novotnyradekcz/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: gator browse [limit]")
	}
	if len(cmd.args) == 1 {
		n, err := strconv.Atoi(cmd.args[0])
		if err != nil || n <= 0 {
			return fmt.Errorf("invalid limit %q: must be a positive integer", cmd.args[0])
		}
		limit = n
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("get posts for user: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found. Follow some feeds and run `gator agg` to collect posts.")
		return nil
	}

	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.Title, post.FeedName)
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", post.PublishedAt.Time.Local().Format("Mon Jan 2, 2006 15:04"))
		}
		fmt.Printf("Link: %s\n", post.Url)
		if post.Description.Valid && post.Description.String != "" {
			fmt.Printf("    %s\n", post.Description.String)
		}
		fmt.Println("=====================================")
	}
	return nil
}
