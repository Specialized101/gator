package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Specialized101/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument\nusage: go run . %s <name>", cmd.name)
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatalf("failed to log in as %s: user does not exist", cmd.args[0])
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("The user %s has been set\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument\nusage: go run . %s <name>", cmd.name)
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		log.Fatalf("error registering new user: %v", err)
	}
	fmt.Printf("user %s was created:\n", cmd.args[0])
	fmt.Printf("  - id: %v\n", user.ID)
	fmt.Printf("  - created_at: %v\n", user.CreatedAt)
	fmt.Printf("  - updated_at: %v\n", user.UpdatedAt)
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		log.Fatal("failed to reset the users table")
	}
	fmt.Println("all records in the users table has been deleted")
	return nil
}
