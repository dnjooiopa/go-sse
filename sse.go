package main

type sseBroker struct {
	users   map[int64][]chan []byte
	actions chan func()
}

func newSSEBroker() *sseBroker {
	s := &sseBroker{
		users:   make(map[int64][]chan []byte),
		actions: make(chan func()),
	}
	go s.run()
	return s
}

func (s *sseBroker) run() {
	for a := range s.actions {
		a()
	}
}

func (s *sseBroker) addUserChan(id int64, ch chan []byte) {
	s.actions <- func() {
		s.users[id] = append(s.users[id], ch)
	}
}

func (s *sseBroker) removeUserChan(id int64, ch chan []byte) {
	go func() {
		for range ch {
		}
	}()

	s.actions <- func() {
		chs := s.users[id]

		i := 0
		for _, c := range chs {
			if c != ch {
				chs[i] = c
				i++
			}
		}
		if i == 0 {
			delete(s.users, id)
		} else {
			s.users[id] = chs[:i]
		}

		close(ch)
	}
}

func (s *sseBroker) sendToUser(id int64, data []byte) {
	s.actions <- func() {
		for _, ch := range s.users[id] {
			ch <- data
		}
	}
}

func (s *sseBroker) userIDs() []int64 {
	var ids []int64
	for id := range s.users {
		ids = append(ids, id)
	}
	return ids
}
