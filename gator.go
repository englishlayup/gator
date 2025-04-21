package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/englishlayup/gator/internal/config"
	"github.com/englishlayup/gator/internal/database"
	"github.com/google/uuid"
)

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	return s.db.DeleteUsers(ctx)
}

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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Expect a single argument, the username.")
	}
	ctx := context.Background()
	name := cmd.args[0]

	_, err := s.db.GetUser(
		ctx,
		sql.NullString{
			String: name,
			Valid:  true,
		},
	)
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

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exist := c.handlers[cmd.name]
	if !exist {
		return fmt.Errorf("Command '%v' not registered", cmd.name)
	}

	return handler(s, cmd)
}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}
