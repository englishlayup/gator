package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func handlerFollowing(s *state, _ command) error {
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, sql.NullString{String: s.cfg.CurrentUserName, Valid: true})
	if err != nil {
		return err
	}
	followingFeeds, err := s.db.GetFeedFollowsForUser(ctx, uuid.NullUUID{UUID: user.ID, Valid: true})

	for _, feed := range followingFeeds {
		fmt.Println(feed.FeedName.String)
	}
	return nil
}
