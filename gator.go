package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/englishlayup/gator/internal/config"
	"github.com/englishlayup/gator/internal/database"
	"github.com/englishlayup/gator/internal/rssfeed"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
		publishDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishDate,
			FeedID:      feed.ID,
		})
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23505" {
					continue
				}
			}
			return err
		}
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

func (c *commands) register(name string, f func(*state, command) error, description string) {
	c.handlers[name] = f
    c.descriptions[name] = description
}

func (c *commands) help() {
    for cmd, desc := range c.descriptions {
        fmt.Println(cmd + ": " + desc)
    }
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
    descriptions map[string]string
}
