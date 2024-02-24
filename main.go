package main

import (
	"log"
	"os"

	"github.com/cbsinteractive/jakewan/distill/cmd"
	"github.com/cbsinteractive/jakewan/distill/production"
)

func init() {
	// Omit date and time.
	log.SetFlags(0)
}

func main() {
	quit := func(err error) {
		log.Fatal(err)
	}
	if err := cmd.Execute(os.Args[1:], production.NewDependencies()); err != nil {
		switch err.Error() {
		case "flag error displayed":
		default:
			quit(err)
		}
	}
}
