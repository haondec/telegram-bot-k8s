package main

import (
	"errors"
	"fmt"
	"os"
	"io/ioutil"
	"time"
	"github.com/gn1k/telegram-bot-k8s/bot"
//	"github.com/go-chat-bot/bot"
)

// Define struct kubebot
type Kubebot struct {
	token    string
	channels map[string]bool
	commands map[string]bool
	roles	 map[string]string
	logname  string
}

// Define constant var will use
const (
	// Declare message announce
	forbiddenCommandMessage  string = "%s - ⚠ Command kubectl %s forbidden\n"
	forbiddenFlagMessage     string = "%s - ⚠ Flag(s) %s forbidden\n"
	forbiddenChannelResponse string = "Sorry @%s, but I'm not allowed to run this command here :zipper_mouth_face:"
	forbiddenCommandResponse string = "Sorry @%s, but I cannot run this command."
	forbiddenFlagResponse    string = "[%s]\nUnknown flag \"%s\".\nCancel task.\n"
	forbiddenFlagResponse_log    string = "Unknown flag: %s."
	forbiddenProjectResponse string = "[%s]\nProject \"%s\" not found.\nCancel task.\n"
	forbiddenProjectResponse_log string = "Project: %s not found."
	// Using
	unAuthorizedUserResponse string = "[%s]\nUnauthorized user: %s.\nCancel task.\n"
	unAuthorizedUserResponse_log string = "Unauthorized user.\n"
	notAllowCommandResponse	 string = "[%s]\n[%s] Not allow to run \"%s\" command.\nPermission denied.\n"
	notAllowCommandResponse_log  string = "%s - Not allow to run: %s command.Permission denied."
	okResponse               string = "[%s]\n%s\n"
	
	// Declare role level
	rolelv3			 string = "projectManager"
	rolelv2			 string = "developer"
	rolelv1			 string = "guest"

	// Format
	timeFM			 string = time.RFC1123Z

	// Deploy help
	deploy_help		 string = `[%s]
Usage: /deploy [OPTION]... [PROJECT NAME] [ENVIROMENT]
Deploy pod, service or deployment on production or other env.
Arguments support.
    -h, --help             show help using
    -s, --show             show list project
    [Project name] [ENV]   /deploy projectA prod`

	// Help bot
	bot_help		 string = `[%s]
@Telegram bot communicate which kubernetes - hcmus
@Support: /help             Get user-id, info command and support
          /deploy           Deploy production, /deploy -h: get more
          /kubectl          All command kubectl support api 1.24

Your id telegram: %s
Contact Sysadmin/ProjectManager to authorize user`
)

// Define var: mapping role <-> user
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
	depcmd = map[string]map[string]bool{
		"proname": map[string]bool{
			"-p": true,
			"--prod": true,
			"--production": true,
			"prod": true,
			"production": true,
		},
	}
)

// Define var: command flag not accep
// No use now
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

//------------------------------------------------------------------------

// Func kubectl [option]... [flag]...
func kubectl(command *bot.Cmd) (msg string, err error) {
	t := time.Now()
	time := t.Format(timeFM)
	userid := command.User.ID
	allow := false
	exist := false

//	fmt.Printf("This is nickname: %s\n", command.User.ID)
//	fmt.Printf("This is id: %s\n", command.User.ID)
//	fmt.Printf("This is realname: %s\n", command.User.RealName)
	
	// Write log recv command
	writeLog(userid, "Receive command kubectl.")

	// Get role of user
	kb_main.roles = rolemap(os.Getenv(telegramRolesLabel))
	rls, exist := kb_main.roles[userid]

	// Checking authorized user
	if !exist {
		writeLog(userid, unAuthorizedUserResponse_log)
		fmt.Printf(unAuthorizedUserResponse, time, userid)
                return fmt.Sprintf(unAuthorizedUserResponse, time, userid), nil
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
		writeLog(userid, fmt.Sprintf(notAllowCommandResponse_log, rls, "kubectl " + command.Args[0]))
		fmt.Printf(notAllowCommandResponse, time, rls, "kubectl " + command.Args[0])
		return fmt.Sprintf(notAllowCommandResponse, time, rls, "kubectl " + command.Args[0]), nil
	}

	return fmt.Sprintf(okResponse, time, output), nil
}

//------------------------------------------------------------------------

// Function deploy
func deploy(command *bot.Cmd) (msg string, err error) {
	t := time.Now()
        time := t.Format(timeFM)
	userid := command.User.ID
	
	// Write log recv command 
        writeLog(userid, "Receive command deploy.")

	// Get role
	kb_main.roles = rolemap(os.Getenv(telegramRolesLabel))
	rls, exist := kb_main.roles[userid]

        // Checking authorized user
        if !exist {
		writeLog(userid, unAuthorizedUserResponse_log)
		fmt.Printf(unAuthorizedUserResponse, time)
                return fmt.Sprintf(unAuthorizedUserResponse, time), nil
        }
	
	// Only /deploy
	if len(command.Args) < 1 {
		// Show help using
		return fmt.Sprintf(deploy_help, time), nil
	}

	// if not Project Manager. Do nothing.
	if rls == rolelv3 {
		writeLog(userid, fmt.Sprintf(notAllowCommandResponse_log, rls, "deploy " + command.Args[0]))
		fmt.Printf(notAllowCommandResponse, time, rls, "deploy " + command.Args[0])
                return fmt.Sprintf(notAllowCommandResponse, time, rls, "deploy " + command.Args[0]), nil
	}
	
	output := ""
	// execute command with deploy
	switch command.Args[0] {
		case "-h", "--help":
			// Show help using
			return fmt.Sprintf(deploy_help, time), nil
		case "-s", "--show":
			// Show list project
			files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
			output = "All project list bellow [Total %d]:\n"
			cnt := 0
			if err != nil {
				output = fmt.Sprintf(output, cnt)
				return fmt.Sprintf(okResponse, time, output), nil
			}
			for _, f := range files {
				if f.IsDir() {
					output += f.Name() + "\n"
					cnt++
				}
			}
			output = fmt.Sprintf(output, cnt)
			return fmt.Sprintf(okResponse, time, output), nil
		default:
			// Unknown flag
			if len(command.Args) < 2 {
				writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[0]))
				fmt.Printf(forbiddenFlagResponse, time, command.Args[0])
				return fmt.Sprintf(forbiddenFlagResponse, time, command.Args[0]), nil
			}
			// Over command
			if len(command.Args) > 2 {
				writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[2]))
				fmt.Printf(forbiddenFlagResponse, time, command.Args[2])
				return fmt.Sprintf(forbiddenFlagResponse, time, command.Args[2]), nil
			}
			
			// Project not found
			proname := command.Args[0]
			check := false
			files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
			if err != nil {
				writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
				fmt.Printf(forbiddenProjectResponse, time, proname)
                                return fmt.Sprintf(forbiddenProjectResponse, time, proname), nil
			}
			for _, f := range files {
				if f.IsDir() && f.Name() == proname {
					check = true
					break
				}
			}

			if check {			// founded
				// This version support only flag env: production or prod
				_, exist := depcmd["proname"][command.Args[1]]
				if exist {
					// Deploy project
					path := os.Getenv(telegramProjectLabel) + proname + "/" + proname + "_prod.yaml"
					kube_command := []string{"create", "-f", path}
					output = execute("kubectl", kube_command...)
					return fmt.Sprintf(okResponse, time, output), nil
				} else {
					writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[1]))
					fmt.Printf(forbiddenFlagResponse, time, command.Args[1])
					return fmt.Sprintf(forbiddenFlagResponse, time, command.Args[1]), nil
				}
			} else {			// not found
				writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
				fmt.Printf(forbiddenProjectResponse, time, proname)
				return fmt.Sprintf(forbiddenProjectResponse, time, proname), nil
			}
	}
}

//------------------------------------------------------------------------

// Function help bot
func info(command *bot.Cmd) (msg string, err error) {
	t := time.Now()
        time := t.Format(timeFM)
        userid := command.User.ID

	// Write log recv command
        writeLog(userid, "Receive command info.")
	return fmt.Sprintf(bot_help, time, userid), nil
}

//------------------------------------------------------------------------

// Init command will use
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
	bot.RegisterCommand(
		"info",
		"Info Telegram integration",
		"",
		info)
}

// Func map file roles of user (file .json) 
func rolemap(fn string) map[string]string {
	claims := getClaims(fn)
	var rm map[string]string
	rm = make(map[string]string)
	for _, p:= range claims {
		rm[p.UserName] = p.Role
	}
	return rm
}
