package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/englishlayup/gator/internal/config"
	"github.com/englishlayup/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	userConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", userConfig.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	gatorState := state{
		dbQueries,
		&userConfig,
	}

	commands := commands{
        handlers: make(map[string]func(*state, command) error),
        descriptions: make(map[string]string),
    }
	commands.register("login", handlerLogin, "Login")
	commands.register("register", handlerRegister, "Register user")
	commands.register("reset", handlerReset, "Reset database")
	commands.register("users", handlerUsers, "List users")
	commands.register("agg", handlerAgg, "Aggregate feeds")
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed), "Add feed to database")
	commands.register("feeds", handlerFeeds, "List available feeds")
	commands.register("follow", middlewareLoggedIn(handlerFollow), "Follow feed")
	commands.register("following", middlewareLoggedIn(handlerFollowing), "List user's following feeds")
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow), "Unfollow feed")
	commands.register("browse", middlewareLoggedIn(handlerBrowse), "Browse posts in following feed")
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: gator [login|register] <username>")
	}
    if args[1] == "help" {
        commands.help()
        return
    }
	command := command{
		name: args[1],
		args: args[2:],
	}
	if err := commands.run(&gatorState, command); err != nil {
		log.Fatal(err)
	}
}
