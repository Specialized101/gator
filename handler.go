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
		return fmt.Errorf("failed to get current user: %w", err)
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
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
		return fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Printf("user %s was created:\n", cmd.args[0])
	fmt.Printf("  - id: %v\n", user.ID)
	fmt.Printf("  - created_at: %v\n", user.CreatedAt)
	fmt.Printf("  - updated_at: %v\n", user.UpdatedAt)
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("failed to delete all users: %w", err)
	}
	fmt.Println("all records in the users table has been deleted")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch all users from db: %w", err)
	}
	for _, u := range users {
		prefix := "* "
		txt := prefix + u.Name

		if u.Name == s.cfg.CurrentUserName {
			txt += " (current)"
		}
		fmt.Println(txt)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", "https://www.wagslane.dev/index.xml", err)
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing arguments\nusage: go run . addFeed <name> <url>")
	}

	_, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	fmt.Printf("RSS Feed '%s' has been added successfully\n", cmd.args[0])

	err = handlerFollow(
		s, command{
			name: "follow",
			args: []string{cmd.args[1]},
		},
		user,
	)
	if err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch all feeds from db: %w", err)
	}
	for i, feed := range feeds {
		u, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to fetch user by id from db: %w", err)
		}
		fmt.Printf("%d. %s:\n", i+1, feed.Name)
		fmt.Printf("  - %s\n", feed.Url)
		fmt.Printf("  - created by %s\n", u.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf(`missing argument. usage: go run . follow "<url>`)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatal("failed to get feed from db")
		return fmt.Errorf("failed to get feed by url from db: %w", err)
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}
	fmt.Printf("User %s just followed '%s'\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows of current user from db: %w", err)
	}
	fmt.Println("Here are the rss feeds you are currently following:")
	for i, ff := range feedFollows {
		i++
		fmt.Printf("	%d. %s\n", i, ff.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing argument\nusage: go run . %s <url>", cmd.name)
	}
	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("failed to delete feed follow: %w", err)
	}
	fmt.Printf("You have unfollowed \"%s\"\n", cmd.args[0])
	return nil
}
