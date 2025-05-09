package main

import "context"

func handlerReset(s *state, _ command) error {
	ctx := context.Background()
	return s.db.DeleteUsers(ctx)
}
