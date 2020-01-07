package command

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Command struct {
	Dir    string
	Args   []string
	Env    map[string]string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type CommandFn func() (int, error)

func (c *Command) Execute() (int, error) {
	commands := map[string]CommandFn{
		"add":    c.cmdAdd,
		"init":   c.cmdInit,
		"commit": c.cmdCommit,
	}

	cmd := c.Args[1]
	f, exists := commands[cmd]
	if !exists {
		fmt.Fprintf(c.Stderr, "jit: '%s' is not a command\n", cmd)
		return 1, nil
	}
	return f()
}

func allEnvVars() map[string]string {
	result := map[string]string{}
	for _, env := range os.Environ() {
		values := strings.Split(env, "=")
		result[values[0]] = values[1]
	}
	return result
}

func New() *Command {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory")
	}

	return &Command{
		Dir:    dir,
		Args:   os.Args,
		Env:    allEnvVars(),
		Stderr: os.Stderr,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}
}
