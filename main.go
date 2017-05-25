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
	notDefinedErrorMessage string = "%s env variable not defined"
)

var (
	kb *Kubebot
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
	return nil
}

func main() {

	if err := validateEnvParams(); err != nil {
		fmt.Printf("Kubebot cannot run. Error: %s\n", err.Error())
		return
	}

	kb = &Kubebot{
		token:    os.Getenv(telegramTokenLabel),
		roles:	  rolemap(os.Getenv(telegramRolesLabel)),
	}

	telegram.Run(kb.token, false)
}
