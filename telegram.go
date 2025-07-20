package pills

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	Z "github.com/rwxrob/bonzai/z"
)

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

func GetTelegramConfig(v *Z.Cmd) (*TelegramConfig, error) {
	var err error
	config := &TelegramConfig{}

	config.BotToken, err = v.Get("telegram.token")
	if err != nil {
		return nil, err
	}
	if config.BotToken == "" {
		return nil, fmt.Errorf("telegram.token is not set")
	}

	config.ChatID, err = v.Get("telegram.chat_id")
	if err != nil {
		return nil, err
	}
	if config.ChatID == "" {
		return nil, fmt.Errorf("telegram.chat_id is not set")
	}

	return config, nil
}

func (config *TelegramConfig) SendMessage(text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.BotToken)

	message := map[string]interface{}{
		"chat_id": config.ChatID,
		"text":    text,
	}

	jsonBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram message: %w", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned non-200 status: %s", resp.Status)
	}

	return nil
}
