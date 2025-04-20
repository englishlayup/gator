package main

import (
	"errors"
	"fmt"

	"github.com/englishlayup/gator/internal/config"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Expect a single argument, the username.")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
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

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}
