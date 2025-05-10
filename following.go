package main

import (
	"context"
	"fmt"
)

func handlerFollowing(s *state, _ command) error {
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	followingFeeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)

	for _, feed := range followingFeeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}
