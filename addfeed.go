package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/englishlayup/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("Expect 2 arguments, the feed name and url.")
	}
	ctx := context.Background()
	name := cmd.args[0]
	url := cmd.args[1]
	currentTime := time.Now()

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return err
	}
	s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	fmt.Printf("Followed %v (%v)\n", feed.Name, feed.Url)
	return nil
}
