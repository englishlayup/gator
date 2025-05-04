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

	commands := commands{make(map[string]func(*state, command) error)}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", handlerAddFeed)
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: gator [login|register] <username>")
	}

	command := command{
		name: args[1],
		args: args[2:],
	}

	if err := commands.run(&gatorState, command); err != nil {
		log.Fatal(err)
	}
}
