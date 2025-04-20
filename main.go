package main

import (
	"log"
	"os"

	"github.com/englishlayup/gator/internal/config"
)


func main() {
    userConfig, err := config.Read()
    if err != nil {
        log.Fatal(err)
    }
    gatorState := state{ &userConfig }

    commands := commands { make(map[string]func(*state, command) error) }

    commands.register("login", handlerLogin)
    args := os.Args
    if len(args) < 2 {
        log.Fatal("Usage: gator login <username>")
    }

    command := command {
        name: "login",
        args: args[2:],
    }

    if err := commands.run(&gatorState, command); err != nil {
        log.Fatal(err)
    }
}
