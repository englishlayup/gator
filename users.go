package main

import (
	"context"
	"fmt"
)

func handlerUsers(s *state, _ command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.String == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.String)
		} else {
			fmt.Printf("* %v\n", user.String)
		}
	}
	return nil
}
