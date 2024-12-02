package httpsrv

import (
	"testing"
)

func TestZeroCall(t *testing.T) {
	strChan := make((chan string), 1)
	s := New(strChan)
	if len(s.sessionStats) != 0 {
		t.Fatalf(`Expected no session stats`)
	}
}

func TestOneCall(t *testing.T) {
	strChan := make((chan string), 1)
	s := New(strChan)
	id := "session1"
	s.incStats(id)
	for _, ws := range s.sessionStats {
		if ws.id == id && ws.sent != 1 {
			t.Fatalf(`Expected session %s sent stat to be 1, but it is %d`, ws.id, ws.sent)
		}
	}
}

func TestSecondCall(t *testing.T) {
	strChan := make((chan string), 1)
	s := New(strChan)
	id := "session1"
	s.incStats(id)
	s.incStats(id)
	for _, ws := range s.sessionStats {
		if ws.id == id && ws.sent != 2 {
			t.Fatalf(`Expected session %s sent stat to be 2, but it is %d`, ws.id, ws.sent)
		}
	}
}

func TestMultipleSessions(t *testing.T) {
	strChan := make((chan string), 1)
	s := New(strChan)
	id1 := "session1"
	id2 := "session2"
	s.incStats(id1)
	s.incStats(id1)
	s.incStats(id2)
	for _, ws := range s.sessionStats {
		if ws.id == id1 && ws.sent != 2 {
			t.Fatalf(`Expected session %s sent stat to be 2, but it is %d`, ws.id, ws.sent)
		}

		if ws.id == id2 && ws.sent != 1 {
			t.Fatalf(`Expected session %s sent stat to be 1, but it is %d`, ws.id, ws.sent)
		}
	}
}
