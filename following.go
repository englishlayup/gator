package main

import (
	"context"
	"fmt"

	"github.com/englishlayup/gator/internal/database"
)

func handlerFollowing(s *state, _ command, user database.User) error {
	ctx := context.Background()
	followingFeeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, feed := range followingFeeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}
