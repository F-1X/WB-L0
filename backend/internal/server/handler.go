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
		http.Error(w, "failed to subscribe subject", http.StatusInternalServerError)
		return
	}
	defer sub.Unsubscribe()

	err = router.stan.NatsConn().PublishRequest("request", replySubj, []byte(id))
	if err != nil {
		http.Error(w, "failed to publish request", http.StatusInternalServerError)
		return
	}
	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		http.Error(w, "timeout exceeded", http.StatusGatewayTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(msg.Data)
}
