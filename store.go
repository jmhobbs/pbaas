package main

import (
	"sync"
	"time"
)

type ProgressDB interface {
	Create(string, string, uint32)
	Update(string, string, uint32) bool
	Get(string) uint32
}

type InMemoryProgressDB struct {
	sync.Mutex
	pbs    map[string]uint32
	tokens map[string]string
	expire map[string]time.Time
}

func (db *InMemoryProgressDB) GC() {
	for {
		time.Sleep(1000)
		now := time.Now()
		db.Lock()
		for key, t := range db.expire {
			if t.Before(now) {
				delete(db.pbs, key)
				delete(db.expire, key)
				delete(db.tokens, key)
			}
		}
		db.Unlock()
	}
}

func NewInMemoryProgressDB() *InMemoryProgressDB {
	db := &InMemoryProgressDB{pbs: make(map[string]uint32), tokens: make(map[string]string), expire: make(map[string]time.Time)}
	go db.GC()
	return db
}

func (db *InMemoryProgressDB) Create(id, token string, progress uint32) {
	db.Lock()
	defer db.Unlock()
	db.pbs[id] = progress
	db.tokens[id] = token
	db.expire[id] = time.Now().Add(time.Minute * 5)
}

func (db *InMemoryProgressDB) Update(id, token string, progress uint32) bool {
	db.Lock()
	defer db.Unlock()
	tkn, ok := db.tokens[id]
	if tkn == token {
		db.pbs[id] = progress
		db.expire[id] = time.Now().Add(time.Minute * 5)
	}
	return ok
}

func (db *InMemoryProgressDB) Get(id string) uint32 {
	db.Lock()
	defer db.Unlock()
	return db.pbs[id]
}
