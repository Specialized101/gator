package main

import (
	"context"
	"fmt"

	"github.com/Specialized101/gator/internal/database"
)

func middlewareLoogedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get current user from db: %w", err)
		}
		err = handler(s, cmd, currentUser)
		if err != nil {
			return err
		}
		return nil
	}
}
