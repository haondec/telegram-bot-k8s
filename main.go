package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-chat-bot/bot/telegram"
)

const (
	telegramTokenLabel        string = "KUBEBOT_TELEGRAM_TOKEN"
	telegramChannelsLabel     string = "KUBEBOT_TELEGRAM_CHANNELS_IDS"
	telegramCommandsLabel     string = "KUBEBOT_TELEGRAM_VALID_COMMANDS"
	notDefinedErrorMessage string = "%s env variable not defined"
)

var (
	kb *Kubebot
)

func validateEnvParams() error {
	if os.Getenv(telegramTokenLabel) == "" {
		return errors.New(fmt.Sprintf(notDefinedErrorMessage, telegramTokenLabel))
	}
	if os.Getenv(telegramChannelsLabel) == "" {
		return errors.New(fmt.Sprintf(notDefinedErrorMessage, telegramChannelsLabel))
	}
	if os.Getenv(telegramCommandsLabel) == "" {
		return errors.New(fmt.Sprintf(notDefinedErrorMessage, telegramCommandsLabel))
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
		channels: stringToMap(os.Getenv(telegramChannelsLabel), " "),
		commands: stringToMap(os.Getenv(telegramCommandsLabel), " "),
	}

	telegram.Run(kb.token)
}
