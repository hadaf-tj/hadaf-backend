package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shb/internal/configs"

	"github.com/stretchr/testify/require"
)

func TestTelegramNotifier_SendAlert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		var payload map[string]string
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)

		require.Equal(t, "123456", payload["chat_id"])
		require.Equal(t, "hello", payload["text"])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))

	defer server.Close()

	n := NewTelegramNotifier(configs.TelegramConfig{
		Token:   "TEST_TOKEN",
		ChatID:  "123456",
		BaseURL: server.URL,
	})

	err := n.SendAlert("hello")

	require.NoError(t, err)
}
