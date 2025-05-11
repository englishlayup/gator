package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/englishlayup/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
        UserID: user.ID,
        Limit: int32(limit),
    })
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println(post.Title.String)
		fmt.Println("===============================================")
		fmt.Println(post.Description.String)
        fmt.Println(post.Url.String)
	}
	return nil
}
