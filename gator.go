package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/englishlayup/gator/internal/config"
	"github.com/englishlayup/gator/internal/database"
	"github.com/englishlayup/gator/internal/rssfeed"
)

func scrapeFeeds(s *state) error {
    feed, err := s.db.GetNextFeedToFetch(context.Background())
    if err != nil {
        return err
    }
    if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
        return err
    }

    rssFeed, err := fetchFeed(context.Background(), feed.Url)
    if err != nil {
        return err
    }
    
    for _, item := range rssFeed.Channel.Item {
        fmt.Println(item.Title)
    }
    return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
    return func(s *state, c command) error {
        user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
        if err != nil {
            return err
        }
        handler(s, c, user)
        return nil
    }
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
