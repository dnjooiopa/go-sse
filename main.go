package main

import (
	"fmt"
	"log"
	"net/http"
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

	clientID := time.Now().UnixNano()

	log.Printf("client connected: %d\n", clientID)

	ch := make(chan string)
	h.sse.addClient(clientID, ch)
	defer h.sse.removeClient(clientID)

	for {
		select {
		case <-r.Context().Done():
			log.Printf("client disconnected: %d\n", clientID)
			return
		case m := <-ch:
			log.Printf("sending message to client: %d\n", clientID)
			fmt.Fprintf(w, "data: %s\n\n", m)
			flusher.Flush()
		}
	}
}

func startWorker(s *sseBroker) {
	go func() {
		for {
			for _, id := range s.clientIDs() {
				s.clients[id] <- fmt.Sprintf("current time is %v", time.Now().Unix())
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {

	sse := &sseBroker{
		clients: make(map[int64]chan string),
	}

	startWorker(sse)

	handler := &appHandler{sse}
	http.HandleFunc("/stream", handler.sseHandler)

	log.Println("server listening on 8888")
	http.ListenAndServe(":8888", nil)
}
