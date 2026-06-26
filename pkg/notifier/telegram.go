package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"shb/internal/configs"
	"time"
)

type TelegramNotifier struct {
	config configs.TelegramConfig
}

func NewTelegramNotifier(cfg configs.TelegramConfig) *TelegramNotifier {
	// Fallback to official Telegram API if BaseURL isn't explicitly provided
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.telegram.org"
	}
	return &TelegramNotifier{
		config: cfg,
	}
}

func (t *TelegramNotifier) SendAlert(message string) error {
	payload := map[string]string{
		"chat_id":    t.config.ChatID,
		"text":       message,
		"parse_mode": "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal telegram payload: %w", err)
	}

	// Clean single-parentheses URL formatting
	url := fmt.Sprintf("%s/bot%s/sendMessage", t.config.BaseURL, t.config.Token)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create telegram request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send telegram request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read telegram response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"telegram api error: status=%d body=%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return nil
}
