package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s.cfg == nil {
		return fmt.Errorf("state's config was nil")
	}
	handler, exists := c.cmds[cmd.name]
	if !exists {
		return fmt.Errorf("command '%s' does not exist", cmd.name)
	}
	handler(s, cmd)
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
