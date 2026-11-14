package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/novotnyradekcz/gator/internal/config"
	"github.com/novotnyradekcz/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	dbURL := cfg.DBURL
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	s := &state{
		cfg: &cfg,
		db:  database.New(dbConn),
	}

	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("agg", handlerAgg)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: not enough arguments provided; a command name is required")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(s, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
