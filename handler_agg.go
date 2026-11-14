package main

import (
	"context"
	"fmt"
)

const feedURL = "https://www.wagslane.dev/index.xml"

func handlerAgg(_ *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("agg does not accept arguments")
	}

	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}
