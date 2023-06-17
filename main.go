package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type appHandler struct {
	sse *sseBroker
}

func (h *appHandler) sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported!", http.StatusInternalServerError)
		return
	}

	userId, err := strconv.ParseInt(r.Header.Get("X-User-Id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	log.Printf("user connected: %d\n", userId)

	ch := make(chan []byte)
	h.sse.addUserChan(userId, ch)
	defer h.sse.removeUserChan(userId, ch)

	for {
		select {
		case <-r.Context().Done():
			log.Printf("client disconnected: %d\n", userId)
			return
		case m := <-ch:
			log.Printf("sending message to user: %d\n", userId)
			fmt.Fprintf(w, "data: %s\n\n", m)
			flusher.Flush()
		}
	}
}

func startWorker(s *sseBroker) {
	go func() {
		for {
			for _, id := range s.userIDs() {
				s.sendToUser(id, []byte(fmt.Sprintf("current time is %v", time.Now().Unix())))
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {

	sse := newSSEBroker()

	startWorker(sse)

	handler := &appHandler{sse}
	http.HandleFunc("/stream", handler.sseHandler)

	log.Println("server listening on 8888")
	http.ListenAndServe(":8888", nil)
}
