package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type config struct {
	width  uint
	height uint
}

func loadConfig() (*config, error) {
	cfg := &config{
		width:  640,
		height: 480,
	}
	args := os.Args
	switch len(args) {
	case 1:
		goto done
	case 2:
		if err := parseSize(cfg, args[1]); err != nil {
			return nil, err
		}
		goto done
	default:
		return nil, errors.New(fmt.Sprint("invalid arguments:", args))
	}
done:
	return cfg, nil
}

func parseSize(cfg *config, arg string) error {
	arg = strings.TrimSpace(arg)
	tokens := strings.Split(arg, ",")
	switch len(tokens) {
	case 1:
		w, err := strconv.Atoi(tokens[0])
		if err != nil {
			return err
		}
		cfg.width = uint(w)
		cfg.height = uint(w) * 3 / 4
	case 2:
		w, err := strconv.Atoi(tokens[0])
		if err != nil {
			return err
		}
		h, err := strconv.Atoi(tokens[0])
		if err != nil {
			return err
		}
		cfg.width = uint(w)
		cfg.height = uint(h)
	default:
		return errors.New(fmt.Sprint("invalid size value: ", tokens))
	}
	return nil
}
