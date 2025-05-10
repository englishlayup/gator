package main

import (
	"errors"
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Expect a single argument, the time between requests.")
	}
    timeBetweenRequests := cmd.args[0]
	intervalDuration, err := time.ParseDuration(timeBetweenRequests)
	if err != nil {
		return nil
	}
    fmt.Println("Collecting feeds every " + timeBetweenRequests)
    ticker := time.NewTicker(intervalDuration)
    for ; ; <-ticker.C {
        scrapeFeeds(s)
    }
}
