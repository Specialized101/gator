package main

import (
	"fmt"
	"log"

	"github.com/Specialized101/gator/internal/config"
)

func main() {
	fmt.Println("Welcome to Gator!")
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}

	cfg.SetUser("specialized")

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}

	fmt.Println(cfg)
}
