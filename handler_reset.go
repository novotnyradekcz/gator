package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, _ command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("reset failed: %w", err)
	}

	fmt.Println("All users have been deleted.")
	return nil
}
