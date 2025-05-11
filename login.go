package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Expect a single argument, the username.")
	}
	ctx := context.Background()
	name := cmd.args[0]
	_, err := s.db.GetUser(ctx, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to login as %v\n", name)
		log.Fatal(err)
	}
	if err := s.cfg.SetUser(name); err != nil {
		return err
	}
	fmt.Printf("Username set to %v", s.cfg.CurrentUserName)
	return nil
}
