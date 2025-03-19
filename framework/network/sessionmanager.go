package network

import (
	"sl.framework.com/async"
	"sl.framework.com/trace"
	"sync"
	"time"
)

type sessioner interface {
	IsAlive() bool
}

type SessionManager struct {
	sync.Mutex
	mapSocks map[string]sessioner
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mapSocks: make(map[string]sessioner),
	}
}

func (s *SessionManager) Start() error {
	async.AsyncRunCoroutine(func() {
		s.checkStatus()
	})
	return nil
}

func (s *SessionManager) AddNewSession(sid string, session sessioner) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.mapSocks[sid]; !ok {
		s.mapSocks[sid] = session
		trace.Info("addnewsession %v success", sid)
		return nil
	} else {
		trace.Error("addnewsession failed, %s existed", sid)
		return nil
	}
}

func (s *SessionManager) getMapSocks() map[string]sessioner {
	s.Lock()
	defer s.Unlock()

	currMapSocks := make(map[string]sessioner)
	for k, v := range s.mapSocks {
		currMapSocks[k] = v
	}

	return currMapSocks
}

func (s *SessionManager) checkStatus() {
	timer := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-timer.C:
			// check if it's alive
			addrs := make([]string, 0, len(s.mapSocks))
			currMapSocks := s.getMapSocks()
			for k, v := range currMapSocks {
				if !v.IsAlive() {
					addrs = append(addrs, k)
				}
			}

			// remove dead connection
			for _, v := range addrs {
				s.removeConn(v)
			}
		}
	}
}

func (s *SessionManager) removeConn(addr string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.mapSocks[addr]; ok {
		delete(s.mapSocks, addr)
		trace.Info("remove connection %v", addr)
	}
}
