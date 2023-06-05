package main

type sseBroker struct {
	clients map[int64]chan string
}

func (s *sseBroker) removeClient(id int64) {
	delete(s.clients, id)
}

func (s *sseBroker) addClient(id int64, ch chan string) {
	s.clients[id] = ch
}

func (s *sseBroker) clientIDs() []int64 {
	ids := []int64{}
	for id := range s.clients {
		ids = append(ids, id)
	}
	return ids
}
