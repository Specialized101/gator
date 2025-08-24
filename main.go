package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Specialized101/gator/internal/config"
)

func main() {
	fmt.Println("Welcome to Gator!")
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}
	var s state
	s.cfg = &cfg
	var c commands
	c.cmds = make(map[string]func(*state, command) error)
	c.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments were provided\n")
	}
	if len(os.Args) < 3 {
		log.Fatal("the username is required to log in\n")
	}
	c.run(&s, command{
		name: os.Args[1],
		args: os.Args[2:],
	})
}
