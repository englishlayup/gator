package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, _ command) error {
	ctx := context.Background()
	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return nil
	}
	fmt.Printf("%v\n", feed)
	return nil
}
