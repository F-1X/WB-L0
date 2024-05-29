package server

import (
	"net/http"
	"time"
)

func (router *Router) orderHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	replySubj := router.stan.NatsConn().NewInbox()
	sub, err := router.stan.NatsConn().SubscribeSync(replySubj)
	if err != nil {
		http.Error(w, "Failed to subscribe to reply subject", http.StatusInternalServerError)
		return
	}
	defer sub.Unsubscribe()

	err = router.stan.NatsConn().PublishRequest("request_http.subject", replySubj, []byte(id))
	if err != nil {
		http.Error(w, "Failed to publish request", http.StatusInternalServerError)
		return
	}
	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		http.Error(w, "Timeout waiting for response", http.StatusGatewayTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(msg.Data)
}
