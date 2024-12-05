package httpsrv

import "log"

var StatsPrint = func(id string, sent int) {
	log.Printf("session %s has received %d messages\n", id, sent)
}

type sessionStats struct {
	id   string
	sent int
}

func (w *sessionStats) print() {
	StatsPrint(w.id, w.sent)
}

func (w *sessionStats) inc() {
	w.sent++
}

func (s *Server) incStats(id string) {
	// Find and increment.
	for _, ws := range s.sessionStats {
		if ws.id == id {
			ws.inc()
			return
		}
	}
	// Not found, add new.
	s.sessionStatsLock.Lock()
	defer s.sessionStatsLock.Unlock()
	s.sessionStats = append(s.sessionStats, &sessionStats{id: id, sent: 1})
}
