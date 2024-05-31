package server

import (
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
)

func (h *Handler) OrderHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	respChan := make(chan *nats.Msg, 1)

	replySubj := h.stan.NatsConn().NewInbox()
	sub, err := h.stan.NatsConn().Subscribe(replySubj, func(msg *nats.Msg) {
		respChan <- msg
	})

	if err != nil {
		http.Error(w, "failed to subscribe subject", http.StatusInternalServerError)
		return
	}
	defer sub.Unsubscribe()

	err = h.stan.NatsConn().PublishRequest("request", replySubj, []byte(id))
	if err != nil {
		http.Error(w, "failed to publish request", http.StatusInternalServerError)
		return
	}

	select {
	case msg := <-respChan:
		log.Println("total time request", time.Since(start))
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(msg.Data)
	case <-time.After(5 * time.Second):
		http.Error(w, "timeout exceeded", http.StatusGatewayTimeout)
		return
	}
}
