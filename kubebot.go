package main

import (
	"errors"
//	"strings"
	"fmt"
	"os"
	"io/ioutil"
	"time"
	"github.com/gn1k/telegram-bot-k8s/bot"
	telegram "github.com/gn1k/telegram-bot-k8s/bot/telegram"
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
	forbiddenCommandMessage		string = "%s - ⚠ Command kubectl %s forbidden\n"
	forbiddenFlagMessage		string = "%s - ⚠ Flag(s) %s forbidden\n"
	forbiddenChannelResponse	string = "Sorry @%s, but I'm not allowed to run this command here :zipper_mouth_face:"
	forbiddenCommandResponse	string = "Sorry @%s, but I cannot run this command."
	forbiddenFlagResponse		string = "- [%s]\nUnknown flag: %s.\nCancel task.\n"
	forbiddenFlagResponse_log	string = "Unknown flag: %s."
	forbiddenProjectResponse	string = "- [%s]\nProject: %s not found.\nCancel task.\n"
	forbiddenProjectResponse_log	string = "Project: %s not found."
	
	// Using
	unAuthorizedUserResponse	string = "- [%s]\nUnauthorized user: %s.\nCancel task.\n"
	unAuthorizedUserResponse_log	string = "Unauthorized user.\n"
	notAllowCommandResponse		string = "- [%s]\n[%s] Not allow to run #%s command.\nPermission denied.\n"
	notAllowCommandResponse_log	string = "Not allow to run: %s command.Permission denied."

	okResponse			string = "- [%s]\n%s\n"

	deploymentResponse_log		string = "Deploy project: %s image:%s, version:%s, env:%s. Completed."
	deploymentResponse		string = "- [%s]\nDeploy project: %s image:%s, version:%s, env:%s.\nCompleted.\n"

	updateResponse_log		string = "Update project: %s image:%s, version:%s, env:%s. Completed."
	updateResponse			string = "- [%s]\nUpdate project: %s image:%s, version:%s, env:%s.\nCompleted.\n"
	
	deleteResponse_log		string = "Delete/Cancel project:%s, env:%s. Completed."
	deleteResponse			string = "- [%s]\nDelete/Cancel project:%s, env:%s.\nCompleted.\n"

	rollbackResponse_log		string = "Rollback project: %s image:%s, version:%s, env:%s. Completed."
	rollbackResponse		string = "- [%s]\nRollback project: %s image:%s, version:%s, env:%s.\nCompleted.\n"

	infoReceiveCommand_log		string = "Receive command %s."
	infoReceiveCommand		string = "- [%s]\nReceive command %s.\n"

	errorConfigFile_log		string = "Project: %s. Error config file: %s"
	errorConfigFile			string = "- [%s]\nProject: %s. Error config file: %s\n"
	
	errorInfoFile_log		string = "Project: %s. Reading info file error."
	errorInfoFile			string = "- [%s]\nProject: %s. Reading info file error.\n"
	
	errorNoState_log		string = "Project: %s. Get #%s state on info file fail - Cancel task."
	errorNoState			string = "- [%s]\nProject: %s. Get #%s state on info file fail - Cancel task.\n"
	
	errorListTag_log		string = "Project: %s. Error fetch all tag: %s."
	errorListTag			string = "- [%s]\nProject: %s. Error fetch all tag: %s.\n"

	errorImageNotFound_log		string = "Project: %s. Image #%s not found."
	errorImageNotFound		string = "- [%s]\nProject: %s. Image #%s not found.\n"
	
	errorTagNotFound_log		string = "Project: %s. Tag #%s not found."
	errorTagNotFound		string = "- [%s]\nProject: %s. Tag #%s not found.\n"
	
	errorSaveInfo_log		string = "Project: %s. Error save info."
	errorSaveInfo			string = "- [%s]\nProject: %s. Error save info.\n"
	errorSaveInfoResponse		string = "Error save info.\n"
	
	infoProjectLocked_log		string = "Project: %s. Locked by another user. Try later."
	infoProjectLocked		string = "- [%s]\nProject: %s. Locked by another user. Try later."

	infoAlreadyLatest_log		string = "Update project: %s. You are already latest/newest. Up to: #%s is the same with current state."
	infoAlreadyLatest           	string = "- [%s]\nUpdate project: %s. You are already latest/newest. Up to: #%s is the same with current state.\n"

	infoAlreadyDeploy_log		string = "Project: %s. Already deploy. Current state exist in info file."
	infoAlreadyDeploy		string = "- [%s]\nProject: %s. Already deploy. Current state exist in info file.\n"

	infoCancelTask_log		string = "Cancel task %s."
	infoCancelTask			string = "- [%s]\nCancel task %s.\n"

	infoCompleteTask_log		string = "Complete task %s."
	infoCompleteTask		string = "- [%s]\nComplete task %s.\n"

	missingFlagResponse_log		string = "Command %s missing flag."
	missingFlagResponse		string = "- [%s]\nCommand %s missing flag.\n"
	
	errorMessage_v1			string = "error:"
	errorMessage_v2			string = "Error from"
	
	// Show flag. Checking config file
	showFlag_v1		string = "%d. %s\n"
	showFlag_v2		string = "%d. %s\n   (Error/Missing config file: %s)\n"

	// Path yaml, script, info
	// dir-project/project-name/peoject-name_env.yaml
	// need validatePath dir-project
	pathYaml		string = "%s%s/%s_%s.yaml"
	pathScript		string = "%s%s/%s_%s.sh"
	pathInfo		string = "%s%s/%s_%s.json"

	// Target script
	targetDeploy		string = "deploy"
	targetDelete		string = "delete"

	// Declare role level
	rolelv3			string = "projectManager"
	rolelv2			string = "developer"
	rolelv1			string = "guest"

	// Format time
	timeFM			string = time.RFC1123Z
	
	// Version
	defaultTag		string = "latest"
	defaultEnv		string = "prod"		//production. Will update "test"
	
	// Deploy help
	deploy_help		 string = `[%s]
Usage: /deploy [OPTION]... [PROJECT NAME] [ENVIRONMENT]
Deploy pod, service or deployment on production or other env.
Arguments support.
	-h, --help		show help using
	-d, --delete		delete deployment
	-s, --show		show list project
	[Project name] [ENV]	/deploy projectA prod`

	// Help bot
	bot_help		 string = `[%s]
@Telegram bot communicate which kubernetes - hcmus
@Support: /help			Get user-id, info command and support
	  /deploy		Deploy production, /deploy -h: get more
	  /kubectl		All command kubectl support api 1.24

Your id telegram: %s
Contact Sysadmin/ProjectManager to authorize user`
	
	// Update help
	update_help		string = `[%s]
Usage: /update [OPTION]... [PROJECT NAME] [ENVIRONMENT]
Arguments support.
	-h, --help		show help using
	-s, --show		show list project
	[Project name] [ENV]	/update projectA prod
				(Default tag:latest)
				/update projectA
				(Default env: production)`

	// Rollback
	rollback_help		string = `[%s]
Usage: /rollback [OPTION]... [PROJECT NAME] [ENVRIRONMENT]
	-h, --help		show help using
	-s, --show		show list project
	[Project name] [ENV]	/rollback projectA prod
				/rollback projectA
				(Default env: production)`
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
		"environment": map[string]bool{
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


// get time
func getTime() string {
	t := time.Now()
	return t.Format(timeFM)
}

//------------------------------------------------------------------------

// Func kubectl [option]... [flag]...
func kubectl(command *bot.Cmd) (msg string, err error) {
	userid := command.User.ID
	allow := false
	exist := false

//	fmt.Printf("This is nickname: %s\n", command.User.ID)
//	fmt.Printf("This is id: %s\n", command.User.ID)
//	fmt.Printf("This is realname: %s\n", command.User.RealName)
	
	// Write log recv command
	this_task := "Kubectl"
	writeLog(userid, fmt.Sprintf(infoReceiveCommand_log, this_task))
	fmt.Printf(infoReceiveCommand, getTime(), this_task)

	// Get role of user
	kb_main.roles = rolemap(os.Getenv(telegramRolesLabel))
	rls, exist := kb_main.roles[userid]

	// Checking authorized user
	if !exist {
		writeLog(userid, unAuthorizedUserResponse_log)
		fmt.Printf(unAuthorizedUserResponse, getTime(), userid)
                return fmt.Sprintf(unAuthorizedUserResponse, getTime(), userid), nil
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
//		fmt.Printf(forbiddenFlagMessage, getTime(), command.Args)
//		return fmt.Sprintf(forbiddenFlagResponse), nil
//	}
//	fmt.Println(command.Args)

	output := ""
	
	if (exist && allow) || !exist {		// Case allow execute command
		output = execute("kubectl", command.Args...)	
	} else {				// Not allow, permission denied
		writeLog(userid, fmt.Sprintf(notAllowCommandResponse_log, rls, "Kubectl " + command.Args[0]))
		fmt.Printf(notAllowCommandResponse, getTime(), rls, "Kubectl " + command.Args[0])
		return fmt.Sprintf(notAllowCommandResponse, getTime(), rls, "Kubectl " + command.Args[0]), nil
	}

	return fmt.Sprintf(okResponse, getTime(), output), nil
}

//------------------------------------------------------------------------

// Function deploy
func deploy(command *bot.Cmd) (msg string, err error) {
	userid := command.User.ID
	
	// Write log recv command
	this_task := "Deploy"
        writeLog(userid, fmt.Sprintf(infoReceiveCommand_log, this_task))
	fmt.Printf(infoReceiveCommand, getTime(), this_task)

	// Get role
	kb_main.roles = rolemap(os.Getenv(telegramRolesLabel))
	rls, exist := kb_main.roles[userid]

        // Checking authorized user
        if !exist {
		writeLog(userid, unAuthorizedUserResponse_log)
		fmt.Printf(unAuthorizedUserResponse, getTime(), userid)
                return fmt.Sprintf(unAuthorizedUserResponse, getTime(), userid), nil
        }
	
	// if not Project Manager. Do nothing.
	if rls != rolelv3 {
		writeLog(userid, fmt.Sprintf(notAllowCommandResponse_log, rls, "Deploy " + command.Args[0]))
		fmt.Printf(notAllowCommandResponse, getTime(), rls, "Deploy " + command.Args[0])
                return fmt.Sprintf(notAllowCommandResponse, getTime(), rls, "Deploy " + command.Args[0]), nil
	}
	
	// Only /deploy
	if len(command.Args) < 1 {
		// Show help using
		return fmt.Sprintf(deploy_help, getTime()), nil
	}
	
	output := ""
	// execute command with deploy
	switch command.Args[0] {
		case "-h", "--help":
			// Show help using
			temp := fmt.Sprintf(infoCompleteTask, getTime(), this_task)
			output = fmt.Sprintf(deploy_help, getTime())
			fmt.Printf(temp)
			return output, nil
		case "-s", "--show":
			// Show list project
			files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
			output = "All project list bellow [Total %d]:\n"
			cnt := 0
			if err != nil {
				output = fmt.Sprintf(output, cnt)
				return fmt.Sprintf(okResponse, getTime(), output), nil
			}

			dir_parent := validatePath(os.Getenv(telegramProjectLabel))
			// This version only support production env
			env := defaultEnv
			check := false
			proname := ""
			for _, f := range files {
				if f.IsDir() {
					// Increase count
					cnt++
					// Checking file yaml, sh
					f_yaml := fmt.Sprintf(pathYaml, dir_parent, f.Name(), f.Name(), env)
					f_script := fmt.Sprintf(pathScript, dir_parent, f.Name(), f.Name(), env)
					f_info := fmt.Sprintf(pathInfo, dir_parent, f.Name(), f.Name(), env)
					lFile := []string{f_yaml, f_script, f_info}
					
					check = check_lock_v2(dir_parent + f.Name())
					proname = f.Name()
					if check {
						proname += " [locked]"
					}
					rs, c := checkConfigFile(lFile, ",")
					
					if c {
						output += fmt.Sprintf(showFlag_v1, cnt, proname)
					} else {
						output += fmt.Sprintf(showFlag_v2, cnt, proname, rs)
					}
				}
			}
			output = fmt.Sprintf(output, cnt)
			temp := fmt.Sprintf(infoCompleteTask, getTime(), this_task)
			output = fmt.Sprintf(okResponse, getTime(), output)
			fmt.Printf(temp)
			return output, nil
		case "-d", "--delete", "-c", "--cancel":
			number := 0
			// Unknown flag
			if len(command.Args) >= 4 {
				number = len(command.Args) - 1
				writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[number]))
				fmt.Printf(forbiddenFlagResponse, getTime(), command.Args[number])
				return fmt.Sprintf(forbiddenFlagResponse, getTime(), command.Args[number]), nil
			}
		
			// Missing flag
			if len(command.Args) != 3 {
				writeLog(userid, fmt.Sprintf(missingFlagResponse_log, "Delete"))
				fmt.Printf(missingFlagResponse, getTime(), "Delete")
				return fmt.Sprintf(missingFlagResponse, getTime(), "Delete"), nil
			}
			
			check := false
			number = len(command.Args) - 1
			env := command.Args[number]

                        _, exist := depcmd["environment"][env]
                        if exist == false {
                        	// Invalid flag
                                writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[number]))
                                fmt.Printf(forbiddenFlagResponse, getTime(), command.Args[number])
                                return fmt.Sprintf(forbiddenFlagResponse, getTime(), command.Args[number]), nil
                        }
			
			// Project not found
			proname := command.Args[1]
			check = false
			files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
			if err != nil {
				writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
				fmt.Printf(forbiddenProjectResponse, getTime, proname)
                                return fmt.Sprintf(forbiddenProjectResponse, getTime, proname), nil
			}
			
			// Find project
			for _, f := range files {
				if f.IsDir() && f.Name() == proname {
					check = true
					break
				}
			}
			
			// Project not found
			if check == false {
				writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
				fmt.Printf(forbiddenProjectResponse, getTime(), proname)
				return fmt.Sprintf(forbiddenProjectResponse, getTime(), proname), nil
			}
		
			// Project locked
			dir_parent := validatePath(os.Getenv(telegramProjectLabel))
			check = check_lock_v2(dir_parent + proname + "/")
			if check == true {
				writeLog(userid, fmt.Sprintf(infoProjectLocked, proname))
				fmt.Printf(infoProjectLocked, getTime(), proname)
				return fmt.Sprintf(infoProjectLocked, getTime(), proname), nil
			}
		
			// Make lock
			make_lock(dir_parent)
			// Unlock after complete function
			defer un_lock(dir_parent)
		
			//############ Delete/Cancel project
			// Now support only env:prod
			env = defaultEnv
			script := fmt.Sprintf(pathScript, dir_parent, proname, proname, env)
			yaml := fmt.Sprintf(pathYaml, dir_parent, proname, proname, env)
			info := fmt.Sprintf(pathInfo, dir_parent, proname, proname, env)
			lFile := []string{yaml, script, info}

			// Check file script, yaml, info
			rs, c := checkConfigFile(lFile, ",")
			if c == false {
				writeLog(userid, fmt.Sprintf(errorConfigFile_log, proname, rs))
				fmt.Printf(errorConfigFile, getTime(), proname, rs)
				return fmt.Sprintf(errorConfigFile, getTime(), proname, rs), nil
			}

			// Read info from file
                        ain, err := getInfo(info)
                        if err != nil {
                                writeLog(userid, fmt.Sprintf(errorInfoFile_log, proname))
                                fmt.Printf(errorInfoFile, getTime(), proname)
				fmt.Println(err.Error())
                                return fmt.Sprintf(errorInfoFile, getTime(), proname), nil
                        }
			
			var image, version string
			in_Current, check = getCurrent(ain)
		
			if check == true {
				// Alway get image from default
				image = in_Current.Name
				version = in_Current.Tag
			} else {
				writeLog(userid, fmt.Sprintf(errorNoState_log, proname, info_TypeCurrent))
                                fmt.Printf(errorNoState, getTime(), proname, info_TypeCurrent)
                                return fmt.Sprintf(errorNoState, getTime(), proname, info_TypeCurrent), nil
			}
			
			//############ Delete
			// With delete no need specific image and version
			pipe_stdin := []string{targetDelete, image, version}
			kube_command := []string{script}
			output = execute_pipe(pipe_stdin, "sh", kube_command...)

			// Update info
			var ain_new []Info
			ain_new = append(ain_new, in_Default)
			err = saveInfo(info, ain_new)
                        if err != nil {
                                writeLog(userid, fmt.Sprintf(errorSaveInfo_log, proname))
                                fmt.Printf(errorSaveInfo, getTime(), proname)
                                bot.SendMessage(telegram.TBot, command.Channel, errorSaveInfoResponse, command.User)
                        }

			writeLog(userid, fmt.Sprintf(deleteResponse_log, proname, "production"))
			temp := fmt.Sprintf(deleteResponse, getTime(), proname, "production")
			fmt.Printf(temp)
			return fmt.Sprintf(okResponse, getTime(), output) + temp, nil
		default:
			check := false
			number := 0
			dir_parent := validatePath(os.Getenv(telegramProjectLabel))
			// Unknown flag
			if len(command.Args) < 2 && len(command.Args) == 3 {
				writeLog(userid, fmt.Sprintf(missingFlagResponse_log, "Deploy"))
				fmt.Printf(missingFlagResponse, getTime(), "Deploy")
				return fmt.Sprintf(missingFlagResponse, getTime(), "Deploy"), nil
			}
			if len(command.Args) > 4 {
				number = 4
				check = true
			}
			
			if check {
				writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[number]))
				fmt.Printf(forbiddenFlagResponse, getTime(), command.Args[number])
				return fmt.Sprintf(forbiddenFlagResponse, getTime(), command.Args[number]), nil
			}
			
			// Project not found
			proname := command.Args[0]
			check = false
			files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
			if err != nil {
				writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
				fmt.Printf(forbiddenProjectResponse, getTime, proname)
                                return fmt.Sprintf(forbiddenProjectResponse, getTime, proname), nil
			}
			
			// Find project
			for _, f := range files {
				if f.IsDir() && f.Name() == proname {
					check = true
					break
				}
			}
			
			// Project not found
			if check == false {
				writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
				fmt.Printf(forbiddenProjectResponse, getTime(), proname)
				return fmt.Sprintf(forbiddenProjectResponse, getTime(), proname), nil
			}
		
			// Project locked
			check = check_lock_v2(dir_parent + proname + "/")
			if check == true {
				writeLog(userid, fmt.Sprintf(infoProjectLocked, proname))
				fmt.Printf(infoProjectLocked, getTime(), proname)
				return fmt.Sprintf(infoProjectLocked, getTime(), proname), nil
			}
		
			// Make lock
			make_lock(dir_parent)
			// Unlock after complete function
			defer un_lock(dir_parent)

			// This version support only flag env: production or prod
			check	= false
			number	= len(command.Args)
			var image, version string
			env	:= command.Args[number - 1]
		
			_, exist := depcmd["environment"][env]
			if exist {
				if number == 4 {
					if command.Args[1] != "--version" && command.Args[1] != "-v" {
						number = 1
						check = true
					}
					version = command.Args[2]
				}
			} else {
				number -= 1
				check = true
			}
			
			// Invalid flag
			if check {
				writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[number]))
				fmt.Printf(forbiddenFlagResponse, getTime(), command.Args[number])
				return fmt.Sprintf(forbiddenFlagResponse, getTime(), command.Args[number]), nil
			}
			
			//############ Deploy project
			// Now support only env:prod
			env = defaultEnv
			script := fmt.Sprintf(pathScript, dir_parent, proname, proname, env)
			yaml := fmt.Sprintf(pathYaml, dir_parent, proname, proname, env)
			info := fmt.Sprintf(pathInfo, dir_parent, proname, proname, env)
			lFile := []string{yaml, script, info}

			// Check file script, yaml, info
			rs, c := checkConfigFile(lFile, ",")
			if c == false {
				writeLog(userid, fmt.Sprintf(errorConfigFile_log, proname, rs))
				fmt.Printf(errorConfigFile, getTime(), proname, rs)
				return fmt.Sprintf(errorConfigFile, getTime(), proname, rs), nil
			}

			// Read info from file
                        ain, err := getInfo(info)
                        if err != nil {
                                writeLog(userid, fmt.Sprintf(errorInfoFile_log, proname))
                                fmt.Printf(errorInfoFile, getTime(), proname)
				fmt.Println(err.Error())
                                return fmt.Sprintf(errorInfoFile, getTime(), proname), nil
                        }
			
			// Get default
			in_Default, check := getDefault(ain)
			if check == true {
				// Alway get image from default
				image = in_Default.Name
				// If Default tag dont specify, value will get by program
				if in_Default.Tag == "" {
					in_Default.Tag = defaultTag
				}
				// Only get default if user dont specify on command /deploy
				if version == "" {
					version = in_Default.Tag
				}
			} else {
				writeLog(userid, fmt.Sprintf(errorNoState_log, proname, info_TypeDefault))
                                fmt.Printf(errorNoState, getTime(), proname, info_TypeDefault)
                                return fmt.Sprintf(errorNoState, getTime(), proname, info_TypeDefault), nil
			}

			// Get current
			_, check = getCurrent(ain)
			if check == true {
				writeLog(userid, fmt.Sprintf(infoAlreadyDeploy_log, proname))
				fmt.Printf(infoAlreadyDeploy, getTime(), proname)
				return fmt.Sprintf(infoAlreadyDeploy, getTime(), proname), nil
			}

			// Handle tag list registry.hub.docker.com
                        ats, err := getAllTags(trueRepo(image))
                        if err != nil {
                                writeLog(userid, fmt.Sprintf(errorListTag_log, proname, image))
                                fmt.Printf(errorListTag, getTime(), proname, image)
				fmt.Println(err.Error())
                                return fmt.Sprintf(errorListTag, getTime(), proname, image), nil
                        }

                        // Checking image
                        if ats[0].Detail == dt_ImageNotFound {
                                writeLog(userid, fmt.Sprintf(errorImageNotFound_log, proname, image))
                                fmt.Printf(errorImageNotFound, getTime(), proname, image)
                                return fmt.Sprintf(errorImageNotFound, getTime(), proname, image), nil
                        }

                        // Cheking tag of image
                        check = findTagName(ats, version)

                        if check == false {
                                version = image + ":" + version
                                writeLog(userid, fmt.Sprintf(errorTagNotFound_log, proname, version))
                                fmt.Printf(errorTagNotFound, getTime(), proname, version)
                                return fmt.Sprintf(errorTagNotFound, getTime(), proname, version), nil
                        }
			
			// Deploy execute
			kube_command := []string{script}
			pipe_stdin := []string{targetDeploy, image, version}
			output = execute_pipe(pipe_stdin, "sh", kube_command...)
			
			// Update info
                        newid := getTagId(ats, version)
                        new_Current := Info{info_TypeCurrent, image, version, newid}
                        ain = applyCurrent(ain, new_Current)
                        ain = applyRollback(ain, new_Current)
                        err = saveInfo(info, ain)
			
                        if err != nil {
                                writeLog(userid, fmt.Sprintf(errorSaveInfo_log, proname))
                                fmt.Printf(errorSaveInfo, getTime(), proname)
                                bot.SendMessage(telegram.TBot, command.Channel, errorSaveInfoResponse, command.User)
                        }

			writeLog(userid, fmt.Sprintf(deploymentResponse_log, proname, image, version, "production"))
			temp := fmt.Sprintf(deploymentResponse, getTime(), proname, image, version, "production")
			fmt.Printf(temp)
			return fmt.Sprintf(okResponse, getTime(), output) + temp, nil
	}
	return "", nil
}

//------------------------------------------------------------------------
func try(command *bot.Cmd) (msg string, err error) {
	bot.SendMessage(telegram.TBot, command.Channel, "Try of command message.\nNice to meet you.", command.User)
	return "", nil
}

//------------------------------------------------------------------------
func update(command *bot.Cmd) (msg string, err error) {
	userid := command.User.ID

	// Write log recv command
	this_task := "Update"
        writeLog(userid, fmt.Sprintf(infoReceiveCommand_log, this_task))
	fmt.Printf(infoReceiveCommand, getTime(), this_task)

        // Get role
        kb_main.roles = rolemap(os.Getenv(telegramRolesLabel))
        rls, exist := kb_main.roles[userid]

        // Checking authorized user
        if !exist {
                writeLog(userid, unAuthorizedUserResponse_log)
                fmt.Printf(unAuthorizedUserResponse, getTime(), userid)
                return fmt.Sprintf(unAuthorizedUserResponse, getTime(), userid), nil
        }

	// if not Project Manager. Do nothing.
        if rls != rolelv3 {
                writeLog(userid, fmt.Sprintf(notAllowCommandResponse_log, rls, "Update " + command.Args[0]))
                fmt.Printf(notAllowCommandResponse, getTime(), rls, "Update " + command.Args[0])
                return fmt.Sprintf(notAllowCommandResponse, getTime(), rls, "Update " + command.Args[0]), nil
        }
	
        // Only /update
        if len(command.Args) < 1 {
                // Show help using
                return fmt.Sprintf(update_help, getTime()), nil
        }

        output := ""
	switch command.Args[0] {
		case "-h", "--help":
                        // Show help using
                        return fmt.Sprintf(update_help, getTime()), nil
                case "-s", "--show":
			// Show list project
			files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
			output = "All project list bellow [Total %d]:\n"
			cnt := 0
			if err != nil {
				output = fmt.Sprintf(output, cnt)
				return fmt.Sprintf(okResponse, getTime(), output), nil
			}

			dir_parent := validatePath(os.Getenv(telegramProjectLabel))
			// This version only support production env
			env := defaultEnv
			check := false
			proname := ""
			for _, f := range files {
				if f.IsDir() {
					// Increase count
					cnt++
					// Checking file yaml, sh
					f_yaml := fmt.Sprintf(pathYaml, dir_parent, f.Name(), f.Name(), env)
					f_script := fmt.Sprintf(pathScript, dir_parent, f.Name(), f.Name(), env)
					f_info := fmt.Sprintf(pathInfo, dir_parent, f.Name(), f.Name(), env)
					lFile := []string{f_yaml, f_script, f_info}
					
					check = check_lock_v2(dir_parent + f.Name())
					proname = f.Name()
					if check {
						proname += " [locked]"
					}
					rs, c := checkConfigFile(lFile, ",")
					
					if c {
						output += fmt.Sprintf(showFlag_v1, cnt, proname)
					} else {
						output += fmt.Sprintf(showFlag_v2, cnt, proname, rs)
					}
				}
			}
			output = fmt.Sprintf(output, cnt)
			temp := fmt.Sprintf(infoCompleteTask, getTime(), this_task)
			output = fmt.Sprintf(okResponse, getTime(), output)
			fmt.Printf(temp)
			return output, nil
		default:
			check := false
			number := len(command.Args) - 1
			if len(command.Args) > 2 {
				check = true
				number = 2
			}
			
			env := defaultEnv

                        if len(command.Args) == 2 {
				number = len(command.Args) - 1
                                _, exist := depcmd["environment"][command.Args[number]]
                                if exist == false {
                                        check = true
                                } else {
                                        env = command.Args[number]
                                }
                        }

                        // Invalid flag
			if check {
                                writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[number]))
                                fmt.Printf(forbiddenFlagResponse, getTime(), command.Args[number])
                                return fmt.Sprintf(forbiddenFlagResponse, getTime(), command.Args[number]), nil
                        }
			
			// Project not found
                        proname := command.Args[0]
                        check = false
                        files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
                        if err != nil {
                                writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
                                fmt.Printf(forbiddenProjectResponse, getTime, proname)
                                return fmt.Sprintf(forbiddenProjectResponse, getTime, proname), nil
                        }
                         
                        // Find project
                        for _, f := range files {
                                if f.IsDir() && f.Name() == proname {
                                        check = true
                                        break
                                }
                        }       
                         
                        // Project not found
                        if check == false {
                                writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
                                fmt.Printf(forbiddenProjectResponse, getTime(), proname)
                                return fmt.Sprintf(forbiddenProjectResponse, getTime(), proname), nil
                        }
		
			// Project locked
			dir_parent := validatePath(os.Getenv(telegramProjectLabel))
			check = check_lock_v2(dir_parent + proname + "/")
			if check == true {
				writeLog(userid, fmt.Sprintf(infoProjectLocked, proname))
				fmt.Printf(infoProjectLocked, getTime(), proname)
				return fmt.Sprintf(infoProjectLocked, getTime(), proname), nil
			}
		
			// Make lock
			make_lock(dir_parent)
			// Unlock after complete function
			defer un_lock(dir_parent)

			// Version: default - latest
                        // This version support only production|prod
			env = defaultEnv
			script := fmt.Sprintf(pathScript, dir_parent, proname, proname, env)
			yaml := fmt.Sprintf(pathYaml, dir_parent, proname, proname, env)
			info := fmt.Sprintf(pathInfo, dir_parent, proname, proname, env)
			kube_command := []string{script}
			lFile := []string{yaml, script, info}

			// Check file script, yaml, info
                        rs, c := checkConfigFile(lFile, ",")
                        if c == false {
                                writeLog(userid, fmt.Sprintf(errorConfigFile_log, proname, rs))
                                fmt.Printf(errorConfigFile, getTime(), proname, rs)
                                return fmt.Sprintf(errorConfigFile, getTime(), proname, rs), nil
                        }

			// Read info from file
			ain, err := getInfo(info)
			if err != nil {
				writeLog(userid, fmt.Sprintf(errorInfoFile_log, proname))
                                fmt.Printf(errorInfoFile, getTime(), proname)
				fmt.Println(err.Error())
                                return fmt.Sprintf(errorInfoFile, getTime(), proname), nil	
			}
			
			var image, version string
			// Get update
			in_Update, check := getUpdate(ain)
			if check == true {
				image = in_Update.Name
				version = in_Update.Tag
			}
			
			// Get current info
			in_Current, check := getCurrent(ain)
			if check == true {
				if image == "" {
					image = in_Current.Name
				}
				if version == "" {
					version = defaultTag
				}
			} else {
				writeLog(userid, fmt.Sprintf(errorNoState_log, proname, info_TypeCurrent))
				fmt.Printf(errorNoState, getTime(), proname, info_TypeCurrent)
                               	return fmt.Sprintf(errorNoState, getTime(), proname, info_TypeCurrent), nil
			}
			
			// Handle tag list registry.hub.docker.com
			ats, err := getAllTags(trueRepo(image))
			if err != nil {
				writeLog(userid, fmt.Sprintf(errorListTag_log, proname, image))
				fmt.Printf(errorListTag, getTime(), proname, image)
				fmt.Println(err.Error())
				return fmt.Sprintf(errorListTag, getTime(), proname, image), nil
			}
			
			// Checking image
			if ats[0].Detail == dt_ImageNotFound {
				writeLog(userid, fmt.Sprintf(errorImageNotFound_log, proname, image))
				fmt.Printf(errorImageNotFound, getTime(), proname, image)
				return fmt.Sprintf(errorImageNotFound, getTime(), proname, image), nil
			}
			
			// Cheking tag of image
			var b1, b2 bool
			b1 = findTagName(ats, version)
			if b1 == false {
				version = image + ":" + version
                                writeLog(userid, fmt.Sprintf(errorTagNotFound_log, proname, version))
                                fmt.Printf(errorTagNotFound, getTime(), proname, version)
                                return fmt.Sprintf(errorTagNotFound, getTime(), proname, version), nil
			}

			// Checking is the same with current
			if image == in_Current.Name {
				if version == in_Current.Tag {
					_, b2 = findTag(ats, in_Current.Tag, in_Current.Id)
					if b2 {
						image = image + ":" + version + ":" + fmt.Sprintf("%d", in_Current.Id)
						writeLog(userid, fmt.Sprintf(infoAlreadyLatest_log, proname, image))
						fmt.Printf(infoAlreadyLatest, getTime(), proname, image)
						return fmt.Sprintf(infoAlreadyLatest, getTime(), proname, image), nil
					}
				}
			}

			//############ Delete
			// With delete no need specific image and version
			pipe_stdin := []string{targetDelete, image, version}
			output = "#Delete output:\n" +  execute_pipe(pipe_stdin, "sh", kube_command...)
			
                        //############ Deploy project
			pipe_stdin = []string{targetDeploy, image, version}
                        output += ("\n#Deploy output:\n" + execute_pipe(pipe_stdin, "sh", kube_command...))
			
			// Update info
			newid := getTagId(ats, version)
			new_Current := Info{info_TypeCurrent, image, version, newid}
			ain = applyCurrent(ain, new_Current)
			ain = applyRollback(ain, in_Current)
			err = saveInfo(info, ain)
			if err != nil {
				writeLog(userid, fmt.Sprintf(errorSaveInfo_log, proname))
				fmt.Printf(errorSaveInfo, getTime(), proname)
				bot.SendMessage(telegram.TBot, command.Channel, errorSaveInfoResponse, command.User)
			}

			// Output return
                        writeLog(userid, fmt.Sprintf(updateResponse_log, proname, image, version, "production"))
			temp := fmt.Sprintf(updateResponse, getTime(), proname, image, version, "production")
                        fmt.Printf(temp)
                        return fmt.Sprintf(okResponse, getTime(), output) + temp, nil
		
	}
	return "", nil
}

//------------------------------------------------------------------------
func rollback(command *bot.Cmd) (msg string, err error) {
	userid := command.User.ID

	// Write log recv command
	this_task := "Rollback"
        writeLog(userid, fmt.Sprintf(infoReceiveCommand_log, this_task))
	fmt.Printf(infoReceiveCommand, getTime(), this_task)

        // Get role
        kb_main.roles = rolemap(os.Getenv(telegramRolesLabel))
        rls, exist := kb_main.roles[userid]

        // Checking authorized user
        if !exist {
                writeLog(userid, unAuthorizedUserResponse_log)
                fmt.Printf(unAuthorizedUserResponse, getTime(), userid)
                return fmt.Sprintf(unAuthorizedUserResponse, getTime(), userid), nil
        }

	// if not Project Manager. Do nothing.
        if rls != rolelv3 {
                writeLog(userid, fmt.Sprintf(notAllowCommandResponse_log, rls, "Rollback " + command.Args[0]))
                fmt.Printf(notAllowCommandResponse, getTime(), rls, "Rollback " + command.Args[0])
                return fmt.Sprintf(notAllowCommandResponse, getTime(), rls, "Rollback " + command.Args[0]), nil
        }
	
        // Only /rollback
        if len(command.Args) < 1 {
                // Show help using
                return fmt.Sprintf(update_help, getTime()), nil
        }

        output := ""
	switch command.Args[0] {
		case "-h", "--help":
                        // Show help using
                        return fmt.Sprintf(rollback_help, getTime()), nil
                case "-s", "--show":
			// Show list project
			files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
			output = "All project list bellow [Total %d]:\n"
			cnt := 0
			if err != nil {
				output = fmt.Sprintf(output, cnt)
				return fmt.Sprintf(okResponse, getTime(), output), nil
			}

			dir_parent := validatePath(os.Getenv(telegramProjectLabel))
			// This version only support production env
			env := defaultEnv
			check := false
			proname := ""
			for _, f := range files {
				if f.IsDir() {
					// Increase count
					cnt++
					// Checking file yaml, sh
					f_yaml := fmt.Sprintf(pathYaml, dir_parent, f.Name(), f.Name(), env)
					f_script := fmt.Sprintf(pathScript, dir_parent, f.Name(), f.Name(), env)
					f_info := fmt.Sprintf(pathInfo, dir_parent, f.Name(), f.Name(), env)
					lFile := []string{f_yaml, f_script, f_info}
					
					check = check_lock_v2(dir_parent + f.Name())
					proname = f.Name()
					if check {
						proname += " [locked]"
					}
					rs, c := checkConfigFile(lFile, ",")
					
					if c {
						output += fmt.Sprintf(showFlag_v1, cnt, proname)
					} else {
						output += fmt.Sprintf(showFlag_v2, cnt, proname, rs)
					}
				}
			}
			output = fmt.Sprintf(output, cnt)
			temp := fmt.Sprintf(infoCompleteTask, getTime(), this_task)
			output = fmt.Sprintf(okResponse, getTime(), output)
			fmt.Printf(temp)
			return output, nil
		default:
			check := false
			number := len(command.Args) - 1
			if len(command.Args) > 2 {
				check = true
				number = 2
			}

                        env := defaultEnv

                        if len(command.Args) == 2 {
                                _, exist := depcmd["environment"][command.Args[number]]
                                if exist == false {
                                        check = true
                                } else {
                                        env = command.Args[number]
                                }       
                        }

			if check {
                                writeLog(userid, fmt.Sprintf(forbiddenFlagResponse_log, command.Args[number]))
                                fmt.Printf(forbiddenFlagResponse, getTime(), command.Args[number])
                                return fmt.Sprintf(forbiddenFlagResponse, getTime(), command.Args[number]), nil
                        }
			
			// Project not found
                        proname := command.Args[0]
                        check = false
                        files, err := ioutil.ReadDir(os.Getenv(telegramProjectLabel))
                        if err != nil {
                                writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
                                fmt.Printf(forbiddenProjectResponse, getTime, proname)
                                return fmt.Sprintf(forbiddenProjectResponse, getTime, proname), nil
                        }
                         
                        // Find project
                        for _, f := range files {
                                if f.IsDir() && f.Name() == proname {
                                        check = true
                                        break
                                }
                        }       
                         
                        // Project not found
                        if check == false {
                                writeLog(userid, fmt.Sprintf(forbiddenProjectResponse_log, proname))
                                fmt.Printf(forbiddenProjectResponse, getTime(), proname)
                                return fmt.Sprintf(forbiddenProjectResponse, getTime(), proname), nil
                        }
		
			// Project locked
			dir_parent := validatePath(os.Getenv(telegramProjectLabel))
			check = check_lock_v2(dir_parent + proname + "/")
			if check == true {
				writeLog(userid, fmt.Sprintf(infoProjectLocked, proname))
				fmt.Printf(infoProjectLocked, getTime(), proname)
				return fmt.Sprintf(infoProjectLocked, getTime(), proname), nil
			}
		
			// Make lock
			make_lock(dir_parent)
			// Unlock after complete function
			defer un_lock(dir_parent)

			// Version: default - latest
                        // This version support only production|prod
			env = defaultEnv
			script := fmt.Sprintf(pathScript, dir_parent, proname, proname, env)
			yaml := fmt.Sprintf(pathYaml, dir_parent, proname, proname, env)
			info := fmt.Sprintf(pathInfo, dir_parent, proname, proname, env)
			kube_command := []string{script}
			lFile := []string{yaml, script, info}

			// Check file script, yaml, info
                        rs, check := checkConfigFile(lFile, ",")
                        if check == false {
                                writeLog(userid, fmt.Sprintf(errorConfigFile_log, proname, rs))
                                fmt.Printf(errorConfigFile, getTime(), proname, rs)
                                return fmt.Sprintf(errorConfigFile, getTime(), proname, rs), nil
                        }

			// Read info from file
			ain, err := getInfo(info)
			if err != nil {
				writeLog(userid, fmt.Sprintf(errorInfoFile_log, proname))
                                fmt.Printf(errorInfoFile, getTime(), proname)
				fmt.Println(err.Error())
                                return fmt.Sprintf(errorInfoFile, getTime(), proname), nil	
			}
			
			// Get rollback info
			in_Rollback, check := getRollback(ain)
			if check == false {
				writeLog(userid, fmt.Sprintf(errorNoState_log, proname, info_TypeRollback))
				fmt.Printf(errorNoState, getTime(), proname, info_TypeRollback)
                               	return fmt.Sprintf(errorNoState, getTime(), proname, info_TypeRollback), nil
			}

			in_Current, check := getCurrent(ain)
			if check == false {
				writeLog(userid, fmt.Sprintf(errorNoState_log, proname, info_TypeCurrent))
                                fmt.Printf(errorNoState, getTime(), proname, info_TypeCurrent)
                                return fmt.Sprintf(errorNoState, getTime(), proname, info_TypeCurrent), nil
			}
			
			image := in_Rollback.Name
			version := in_Rollback.Tag
			
			// Handle tag list registry.hub.docker.com
			ats, err := getAllTags(trueRepo(image))
			if err != nil {
				writeLog(userid, fmt.Sprintf(errorListTag_log, proname, image))
				fmt.Printf(errorListTag, getTime(), proname, image)
				fmt.Println(err.Error())
				return fmt.Sprintf(errorListTag, getTime(), proname, image), nil
			}
			
			// Checking image
			if ats[0].Detail == dt_ImageNotFound {
				writeLog(userid, fmt.Sprintf(errorImageNotFound_log, proname, image))
				fmt.Printf(errorImageNotFound, getTime(), proname, image)
				return fmt.Sprintf(errorImageNotFound, getTime(), proname, image), nil
			}
			
			// Cheking tag of image
			check = findTagName(ats, version)
			if check == false {
				version = image + ":" + version
                                writeLog(userid, fmt.Sprintf(errorTagNotFound_log, proname, version))
                                fmt.Printf(errorTagNotFound, getTime(), proname, version)
                                return fmt.Sprintf(errorTagNotFound, getTime(), proname, version), nil
			}

			//############ Delete
			// With delete no need specific image and version
			pipe_stdin := []string{targetDelete, image, version}
			output = "#Delete output:\n" +  execute_pipe(pipe_stdin, "sh", kube_command...)
			
                        //############ Deploy project
			pipe_stdin = []string{targetDeploy, image, version}
                        output += ("\n#Deploy output:\n" + execute_pipe(pipe_stdin, "sh", kube_command...))
			
			// Update info
			newid := getTagId(ats, version)
			new_Current := Info{info_TypeCurrent, image, version, newid}
			ain = applyCurrent(ain, new_Current)
			ain = applyRollback(ain, in_Current)
			err = saveInfo(info, ain)
			if err != nil {
				writeLog(userid, fmt.Sprintf(errorSaveInfo_log, proname))
				fmt.Printf(errorSaveInfo, getTime(), proname)
				bot.SendMessage(telegram.TBot, command.Channel, errorSaveInfoResponse, command.User)
			}

			// Output return
                        writeLog(userid, fmt.Sprintf(rollbackResponse_log, proname, image, version, "production"))
			temp := fmt.Sprintf(rollbackResponse, getTime(), proname, image, version, "production")
                        fmt.Printf(temp)
                        return fmt.Sprintf(okResponse, getTime(), output) + temp, nil
		
	}
	return "", nil
}

//------------------------------------------------------------------------

// Function help bot
func info(command *bot.Cmd) (msg string, err error) {
        userid := command.User.ID

	// Write log recv command
        writeLog(userid, "Receive command info.")
	return fmt.Sprintf(bot_help, getTime(), userid), nil
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
	bot.RegisterCommand(
		"update",
		"Update Telegram integration",
		"",
		update)
	bot.RegisterCommand(
		"rollback",
		"Rollback Telegram integration",
		"",
		rollback)
}
