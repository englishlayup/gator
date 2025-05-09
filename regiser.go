package main

import (
	"context"
	"database/sql"
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
			CreatedAt: sql.NullTime{Time: currentTime, Valid: true},
			UpdatedAt: sql.NullTime{Time: currentTime, Valid: true},
			Name:      sql.NullString{String: name, Valid: true},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User %v created", user.Name.String)

	if err := s.cfg.SetUser(user.Name.String); err != nil {
		return err
	}
	return nil
}
