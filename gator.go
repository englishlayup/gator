package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/englishlayup/gator/internal/config"
	"github.com/englishlayup/gator/internal/database"
	"github.com/englishlayup/gator/internal/rssfeed"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("Expect 2 arguments, the feed name and url.")
	}
	ctx := context.Background()
	name := cmd.args[0]
	url := cmd.args[1]
	currentTime := time.Now()
	user, err := s.db.GetUser(ctx, sql.NullString{String: s.cfg.CurrentUserName, Valid: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't retrieve user %v while adding feed\n", user.Name)
		return err
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: currentTime, Valid: true},
		UpdatedAt: sql.NullTime{Time: currentTime, Valid: true},
		Name:      sql.NullString{String: name, Valid: true},
		Url:       sql.NullString{String: url, Valid: true},
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	}
	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", feed)
	return nil
}

func handlerAgg(s *state, _ command) error {
	ctx := context.Background()
	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return nil
	}
	fmt.Printf("%v\n", feed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*rssfeed.RSSFeed, error) {
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	var rssFeed rssfeed.RSSFeed
	if err := xml.Unmarshal(body, &rssFeed); err != nil {
		return nil, err
	}
	unescapeRSSFeed(&rssFeed)
	return &rssFeed, nil
}

func unescapeRSSFeed(feed *rssfeed.RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
}

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

func handlerReset(s *state, _ command) error {
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
