package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/englishlayup/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("Expect 2 arguments, the feed name and url.")
	}
	ctx := context.Background()
	name := cmd.args[0]
	url := cmd.args[1]
	currentTime := time.Now()
	user, err := s.db.GetUser(ctx, sql.NullString{String: s.cfg.CurrentUserName, Valid: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't retrieve user %v while adding feed\n", user.Name)
		return err
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: currentTime, Valid: true},
		UpdatedAt: sql.NullTime{Time: currentTime, Valid: true},
		Name:      sql.NullString{String: name, Valid: true},
		Url:       sql.NullString{String: url, Valid: true},
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	}
	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return err
	}
	s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: currentTime, Valid: true},
		UpdatedAt: sql.NullTime{Time: currentTime, Valid: true},
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
	})

	fmt.Printf("Followed %v (%v)\n", feed.Name.String, feed.Url.String)
	return nil
}
