package main

import (
	"strings"

	"os"
	"os/signal"
	"syscall"

	"github.com/apoloval/simavionics"
	"github.com/apoloval/simavionics/a320"
	"github.com/chzyer/readline"
	"github.com/op/go-logging"
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
		case "log":
			c.Log(args)
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

func (c *CLI) Log(args []string) {
	var level = logging.NOTICE
	switch {
	case argsMatch(args, "notice"):
		level = logging.NOTICE
	case argsMatch(args, "info"):
		level = logging.INFO
	case argsMatch(args, "warning"):
		level = logging.WARNING
	case argsMatch(args, "error"):
		level = logging.ERROR
	case len(args) == 0:
	default:
		printSyntaxError("log [notice|info|warning|error]")
	}
	println("You are entering in the log viewing mode.")
	println("Press CTRL+C to go back to the simulation console.")
	println()
	simavionics.EnableLoggingLevel(level)
	waitForStopSignal()
	simavionics.DisableLogging()
	println()
	println("Exiting log viewing mode and going back to simulation console.")
	println()
}

func (c *CLI) Apu(args []string) {
	switch {
	case argsMatch(args, "master", "on"):
		simavionics.PublishEvent(c.bus, a320.ApuActionMasterSwOn, true)
	case argsMatch(args, "master", "off"):
		simavionics.PublishEvent(c.bus, a320.ApuActionMasterSwOn, false)
	default:
		printSyntaxError("apu master (on|off)")
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

func waitForStopSignal() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
}
