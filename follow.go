package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/englishlayup/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Expect 1 argument, the url.")
	}
	url := cmd.args[0]
	currentTime := time.Now()
	ctx := context.Background()

	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feedFollowRecord, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(feedFollowRecord.UserName)
	fmt.Println(feedFollowRecord.FeedName)
	return nil
}
