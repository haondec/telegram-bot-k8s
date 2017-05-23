package main

import (
	"errors"
	"fmt"
	"time"
	"github.com/gn1k/telegram-dev/bot"
//	"github.com/go-chat-bot/bot"
)

type Kubebot struct {
	token    string
	channels map[string]bool
	commands map[string]bool
	roles	 map[string]string
}

const (
	// Declare message announce
	forbiddenCommandMessage  string = "%s - ⚠ Command kubectl %s forbidden\n"
	forbiddenFlagMessage     string = "%s - ⚠ Flag(s) %s forbidden\n"
	forbiddenChannelResponse string = "Sorry @%s, but I'm not allowed to run this command here :zipper_mouth_face:"
	forbiddenCommandResponse string = "Sorry @%s, but I cannot run this command."
	forbiddenFlagResponse    string = "Sorry @%s, but I'm not allowed to run one of your flags."
	
	// Using
	unAuthorizedUserResponse string = "[%s]\nUnauthorized user\n"
	notAllowCommandResponse	 string = "[%s]\n[%s] Not allow to run \"%s\" command.\nPermission denied.\n"
	okResponse               string = "[%s]\n%s\n"
	
	// Declare role level
	rolelv3			 string = "projectManager"
	rolelv2			 string = "developer"
	rolelv1			 string = "guest"
)

var (
	rolecmd = map[string]map[string]bool{
		"create": map[string]bool{
			"developer":	false,
			"guest":	false,
		},
		"delete": map[string]bool{
			"developer":	false,
			"guest":	false,
		},
		"run": map[string]bool{
			"developer":	false,
			"guest":	false,
		},
		"exec": map[string]bool{
                        "developer":    false,
                        "guest":        false,
                },
		"scale": map[string]bool{
                        "developer":    false,
                        "guest":        false,
                },
		"apply": map[string]bool{
                        "developer":    false,
                        "guest":        false,
                },
	}
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
	t := time.Now()
	time := t.Format(time.RFC3339)
	userid := command.User.ID
	allow := false
	exist := false

//	fmt.Printf("This is nickname: %s\n", command.User.ID)
//	fmt.Printf("This is id: %s\n", command.User.ID)
//	fmt.Printf("This is realname: %s\n", command.User.RealName)
	
	// Get role of user
	rls, exist := kb.roles[userid]

	// Checking authorized user
	if !exist {
		fmt.Printf(unAuthorizedUserResponse, time)
                return fmt.Sprintf(unAuthorizedUserResponse, time), nil
	} else {
		exist = false
	}
	
	// Checking role command
	if rls  == rolelv3 {			// Project manager
		allow = true
		exist = true
	} else if rls == rolelv2 {		// Developer
		allow, exist = rolecmd[command.Args[0]]["developer"]
	} else if rls == rolelv1 {		// Guest
		allow, exist = rolecmd[command.Args[0]]["guest"]
	} else {				// Unknow role defined
		allow = false
		exist = true
	}
	
//	if err := validateFlags(command.Args...); err != nil {
//		fmt.Printf(forbiddenFlagMessage, time, command.Args)
//		return fmt.Sprintf(forbiddenFlagResponse), nil
//	}
//	fmt.Println(command.Args)

	output := ""
	
	if (exist && allow) || !exist {		// Case allow execute command
		output = execute("kubectl", command.Args...)	
	} else {				// Not allow, permission denied
		fmt.Printf(notAllowCommandResponse, time, rls, "kubectl " + command.Args[0])
		return fmt.Sprintf(notAllowCommandResponse, time, rls, "kubectl " + command.Args[0]), nil
	}

	return fmt.Sprintf(okResponse, time, output), nil
}

func deploy(command *bot.Cmd) (msg string, err error) {
	t := time.Now()
        time := t.Format(time.RFC3339)
	userid := command.User.ID
	rls, exist := kb.roles[userid]

        // Checking authorized user
        if !exist {
                fmt.Printf(unAuthorizedUserResponse, time)
                return fmt.Sprintf(unAuthorizedUserResponse, time), nil
        }
	
	// if not Project Manager. Do nothing.
	if rls != rolelv3 {
		fmt.Printf(notAllowCommandResponse, time, rls, "kubectl " + command.Args[0])
                return fmt.Sprintf(notAllowCommandResponse, time, rls, "kubectl " + command.Args[0]), nil
	}
	
	// execute command with deploy
	output := execute("whoami", command.Args...)

	return fmt.Sprintf(okResponse, time, output), nil
}

func init() {
	bot.RegisterCommand(
		"kubectl",
		"Kubectl Telegram integration",
		"",
		kubectl)

	bot.RegisterCommand(
		"deploy",
		"Deploy Telegram integration",
		"",
		deploy)
}

func rolemap(fn string) map[string]string {
	claims := getClaims(fn)
	var rm map[string]string
	rm = make(map[string]string)
	for _, p:= range claims {
		rm[p.UserName] = p.Role
	}
	return rm
}
