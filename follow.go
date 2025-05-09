package main

import (
	"context"
	"database/sql"
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

	feed, err := s.db.GetFeedByUrl(ctx, sql.NullString{String: url, Valid: true})
	if err != nil {
		return err
	}
	user, err := s.db.GetUser(ctx, sql.NullString{String: s.cfg.CurrentUserName, Valid: true})
	if err != nil {
		return err
	}
	feedFollowRecord, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: currentTime, Valid: true},
		UpdatedAt: sql.NullTime{Time: currentTime, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		return err
	}
	fmt.Println(feedFollowRecord.UserName.String)
	fmt.Println(feedFollowRecord.FeedName.String)
	return nil
}
