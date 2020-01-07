package main

import (
	"log"
	"os"

	"github.com/tpbowden/jit/command"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("No command")
	}

	cmd := command.New()
	status, err := cmd.Execute()
	if err != nil {
		panic(err)
	}
	os.Exit(status)
}
