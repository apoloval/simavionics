package main

import (
	"strings"

	"log"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320/apu"
	"github.com/chzyer/readline"
)

type CLI struct {
	lines *readline.Instance
	bus   simavionics.EventBus
}

func NewCLI(bus simavionics.EventBus) (*CLI, error) {
	lines, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return nil, err
	}

	log.SetOutput(lines.Stderr())

	return &CLI{
		lines: lines,
		bus:   bus,
	}, nil
}

func (c *CLI) Run() {
	println("SimAvionics - A320 Systems Simulation")
	println("Copyright (C) 2018 Alvaro Polo")
	println("")
	for {
		line, err := c.lines.Readline()
		if err != nil {
			goto exit
		}

		tokens := tokenize(line)
		if len(tokens) == 0 {
			continue
		}

		cmd := tokens[0]
		args := tokens[1:]
		switch cmd {
		case "apu":
			c.Apu(args)
		case "pub":
			c.Pub(args)
		case "exit", "quit":
			goto exit
		default:
			println("Invalid command", line)
		}
	}
exit:
	println("Exiting...")
	return
}

func (c *CLI) Pub(args []string) {
	if len(args) != 2 {
		printSyntaxError("pub <event name> <value>")
		return
	}

	event := simavionics.EventName(args[0])
	value := parseValue(args[1])
	if value == nil {
		println("Invalid event value:", args[1])
		return
	}
	simavionics.PublishEvent(c.bus, event, value)
}

func (c *CLI) Apu(args []string) {
	switch {
	case argsMatch(args, "master", "on"):
		simavionics.PublishEvent(c.bus, apu.EventMasterSwitch, true)
	case argsMatch(args, "master", "off"):
		simavionics.PublishEvent(c.bus, apu.EventMasterSwitch, false)
	case argsMatch(args, "start"):
		simavionics.PublishEvent(c.bus, apu.EventStartButton, true)
	case argsMatch(args, "bleed", "on"):
		simavionics.PublishEvent(c.bus, apu.EventBleedSwitch, true)
	case argsMatch(args, "bleed", "off"):
		simavionics.PublishEvent(c.bus, apu.EventBleedSwitch, false)
	default:
		printSyntaxError(
			"apu master (on|off)",
			"apu start",
			"apu bleed (on|off)",
		)
	}
}

func tokenize(line string) []string {
	raw := strings.Split(line, " ")
	var tokens []string
	for _, r := range raw {
		t := strings.TrimSpace(r)
		if len(t) > 0 {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func parseValue(v string) interface{} {
	v = strings.TrimSpace(v)
	switch v {
	case "true", "TRUE":
		return true
	case "false", "FALSE":
		return false
	}
	return nil
}

func argsMatch(args []string, expected ...string) bool {
	if len(args) != len(expected) {
		return false
	}
	for i, e := range expected {
		if args[i] != e {
			return false
		}
	}
	return true
}

func printSyntaxError(usages ...string) {
	println("Syntax error, expected:")
	for _, u := range usages {
		println("   ", u)
	}
}
