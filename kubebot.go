package main

import (
	"errors"
	"fmt"
//	"time"

	"github.com/go-chat-bot/bot"
)

type Kubebot struct {
	token    string
	channels map[string]bool
	commands map[string]bool
}

const (
	forbiddenCommandMessage  string = "%s - ⚠ Command kubectl %s forbidden\n"
	forbiddenFlagMessage     string = "%s - ⚠ Flag(s) %s forbidden\n"
	forbiddenChannelResponse string = "Sorry @%s, but I'm not allowed to run this command here :zipper_mouth_face:"
	forbiddenCommandResponse string = "Sorry @%s, but I cannot run this command."
	forbiddenFlagResponse    string = "Sorry @%s, but I'm not allowed to run one of your flags."
	okResponse               string = "\n@%s\n"
)

var (
	ignored = map[string]map[string]bool{
		"get": map[string]bool{
			"-f":           true,
			"--filename":   true,
			"-w":           true,
			"--watch":      true,
			"--watch-only": true,
		},
		"describe": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"create": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"replace": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"patch": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"delete": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"edit": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"apply": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"logs": map[string]bool{
			"-f":       true,
			"--follow": true,
		},
		"rolling-update": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"scale": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"attach": map[string]bool{
			"-i":      true,
			"--stdin": true,
			"-t":      true,
			"--tty":   true,
		},
		"exec": map[string]bool{
			"-i":      true,
			"--stdin": true,
			"-t":      true,
			"--tty":   true,
		},
		"run": map[string]bool{
			"--leave-stdin-open": true,
			"-i":                 true,
			"--stdin":            true,
			"--tty":              true,
		},
		"expose": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"autoscale": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"label": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"annotate": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"convert": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
	}
)

func validateFlags(arguments ...string) error {
	if len(arguments) <= 1 {
		return nil
	}

	for i := 1; i < len(arguments); i++ {
		if ignored[arguments[0]][arguments[i]] {
			return errors.New(fmt.Sprintf("Error: %s is an invalid flag", arguments[i]))
		}

	}

	return nil
}

func kubectl(command *bot.Cmd) (msg string, err error) {
//	t := time.Now()
//	time := t.Format(time.RFC3339)

//	if err := validateFlags(command.Args...); err != nil {
//		fmt.Printf(forbiddenFlagMessage, time, command.Args)
//		return fmt.Sprintf(forbiddenFlagResponse), nil
//	}

	output := execute("kubectl", command.Args...)

	return fmt.Sprintf(okResponse, output), nil
}

func init() {
	bot.RegisterCommand(
		"kubectl",
		"Kubectl Telegram integration",
		"",
		kubectl)
}
