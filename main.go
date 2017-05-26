package main

import (
	"errors"
	"fmt"
	"os"
	"github.com/gn1k/telegram-dev/bot/telegram"
//	"github.com/go-chat-bot/bot/telegram"
)

const (
	telegramTokenLabel        string = "KUBEBOT_TELEGRAM_TOKEN"
	telegramRolesLabel	  string = "KUBEBOT_TELEGRAM_ROLES_FILE"
	telegramProjectLabel	  string = "KUBEBOT_TELEGRAM_PROJECT_DIR"
	telegramLogLabel	  string = "KUBEBOT_TELEGRAM_LOG_DIR"
	notDefinedErrorMessage string = "%s env variable not defined"
)

var (
	kb_main *Kubebot
	lcfg_main *LogConfig
)

func validateEnvParams() error {
	if os.Getenv(telegramTokenLabel) == "" {
		return errors.New(fmt.Sprintf(notDefinedErrorMessage, telegramTokenLabel))
	}
	if os.Getenv(telegramRolesLabel) == "" {
		return errors.New(fmt.Sprintf(notDefinedErrorMessage, telegramRolesLabel))
	}
	if os.Getenv(telegramProjectLabel) == "" {
		return errors.New(fmt.Sprintf(notDefinedErrorMessage, telegramProjectLabel))
	}
	if os.Getenv(telegramLogLabel) == "" {
		return errors.New(fmt.Sprintf(notDefinedErrorMessage, telegramLogLabel))
	}
	return nil
}

// Validate path
func validatePath(str string) string {
	if str != "" {
		if str[len(str) - 1] != '/' {
			str += "/"
		}
	}
	return str
}

func main() {

	if err := validateEnvParams(); err != nil {
		fmt.Printf("Kubebot cannot run. Error: %s\n", err.Error())
		return
	}
	
	kb_main = &Kubebot{
		token:    os.Getenv(telegramTokenLabel),
		roles:	  rolemap(os.Getenv(telegramRolesLabel)),
	}

	lcfg_main = &LogConfig{
		log_count:	0,
		log_count_max:	1000,
		log_sub:	0,

		log_file_name:	initLogName(),
	}

	telegram.Run(kb_main.token, false)
}
