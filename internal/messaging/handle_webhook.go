package messaging

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func (m *messaging) Callback(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	lineSignature := r.Header.Get("X-Line-Signature")
	if !webhook.ValidateSignature(m.config.LineChannelSecret, lineSignature, body) {
		slog.Error("invalid line signature", slog.String("signature", lineSignature))
		w.WriteHeader(http.StatusOK)
		return
	}

	var cb webhook.CallbackRequest
	if err := json.Unmarshal(body, &cb); err != nil {
		slog.Error("failed to unmarshal request body", slog.Any("error", err))
		w.WriteHeader(http.StatusOK)
		return
	}

	// TODO: reply message
	slog.Info("callback", slog.Any("response", cb))

	w.WriteHeader(http.StatusOK)
}
