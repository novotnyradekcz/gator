package main

import (
	"context"
	"fmt"

	"github.com/novotnyradekcz/gator/internal/database"
)

// middlewareLoggedIn wraps a handler that requires a logged-in user. It looks
// up the current user and passes it to the handler, returning an error if the
// user cannot be found.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("get current user: %w", err)
		}

		return handler(s, cmd, user)
	}
}
