package main

import (
	"context"
	"log"
	"time"

	"github.com/marceljaworski/cli-spinner/spinner"
)

func main() {
	s := spinner.New(spinner.Config{})

	ctx, cancel := context.WithCancel(context.Background())

	log.Println("Starting the spinner")
	s.Start(ctx)
	time.Sleep(time.Second * 5)
	cancel()
	log.Println("Spinner stopped")
}
