package main

import (
	"fmt"
	"log"

	"github.com/novotnyradekcz/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.SetUser("novotnyradekcz"); err != nil {
		log.Fatal(err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)
}
