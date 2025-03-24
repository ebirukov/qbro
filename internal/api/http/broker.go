package http

import (
	"context"
	"ebirukov/qbro/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const queueID = "queueID"

func (h *BrokerHandler) Put(w http.ResponseWriter, r *http.Request) {
	qID := r.PathValue(queueID)
	if qID == "" {
		http.Error(w, "queue param required", http.StatusBadRequest)

		return
	}

	var msg Message

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, fmt.Errorf("can't read message: %w", err).Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	if err = h.broker.Put(ctx, model.QueueID(qID), model.Message(msg.Data)); err != nil {
		http.Error(w, err.Error(), translateError(err))

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BrokerHandler) Get(w http.ResponseWriter, r *http.Request) {
	qID := r.PathValue(queueID)
	if qID == "" {
		http.Error(w, "queue param required", http.StatusBadRequest)

		return
	}
	ctx := r.Context()
	timeoutParam := r.URL.Query().Get("timeout")

	to := h.cfg.DefaultTimeout

	if timeoutParam != "" {
		var err error

		to, err = time.ParseDuration(timeoutParam+"s")
		if err != nil {
			err := fmt.Errorf("can't parse timeout param: %w", err)
	
			http.Error(w, err.Error(), http.StatusBadRequest)
	
			return
		}
	}

	if to > 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, to)
		defer cancel()
	}

	msg, err := h.broker.Get(ctx, model.QueueID(qID))
	if err != nil {
		http.Error(w, err.Error(), translateError(err))

		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(Message{msg.String()})
	if err != nil {
		http.Error(w, fmt.Errorf("can't parse timeout param: %w", err).Error(), http.StatusBadRequest)

		return
	}

}