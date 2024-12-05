package httpsrv

import (
	"goapp/internal/pkg/watcher"
)

func (s *Server) addWatcher(w *watcher.Watcher) {
	s.watchersLock.Lock()
	defer s.watchersLock.Unlock()
	s.watchers[w.GetWatcherId()] = w
}

func (s *Server) removeWatcher(w *watcher.Watcher) {
	s.watchersLock.Lock()
	defer s.watchersLock.Unlock()
	// Print satistics before removing watcher.
	for i := range s.sessionStats {
		if s.sessionStats[i].id == w.GetWatcherId() {
			s.sessionStats[i].print()

			// Remove unneeded stat
			s.sessionStatsLock.Lock()
			sessionStatsLength := len(s.sessionStats)
			s.sessionStats[i] = s.sessionStats[sessionStatsLength-1] // Copy last element to index i.
			s.sessionStats = s.sessionStats[:sessionStatsLength-1]   // Truncate slice.

			// s.sessionStats = append(s.sessionStats[:i], s.sessionStats[i+1:]...)
			s.sessionStatsLock.Unlock()
			break
		}
	}
	// Remove watcher.
	delete(s.watchers, w.GetWatcherId())
}

func (s *Server) notifyWatchers(str string) {
	s.watchersLock.RLock()
	defer s.watchersLock.RUnlock()

	// Send message to all watchers and increment stats.
	for id := range s.watchers {
		s.watchers[id].Send(str)
		s.incStats(id)
	}
}
