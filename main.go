package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Specialized101/gator/internal/config"
	"github.com/Specialized101/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Welcome to Gator!")
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}
	var s state
	s.cfg = &cfg

	db, err := sql.Open("postgres", s.cfg.DbUrl)
	if err != nil {
		log.Fatalf("error opening the database: %v", err)
	}
	s.db = database.New(db)

	var c commands
	c.cmds = make(map[string]func(*state, command) error)
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments were provided\nusage: go run . <command> [arguments]")
	}
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	c.run(&s, command{
		name: os.Args[1],
		args: args,
	})

}
