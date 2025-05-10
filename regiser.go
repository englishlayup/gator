package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/englishlayup/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Expect a single argument, the username.")
	}

	name := cmd.args[0]
	currentTime := time.Now()
	ctx := context.Background()

	user, err := s.db.CreateUser(
		ctx,
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			Name:      name,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User %v created", user.Name)

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}
	return nil
}
