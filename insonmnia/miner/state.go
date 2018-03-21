package miner

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/pborman/uuid"
	"github.com/sonm-io/core/proto"
)

const stateKey = "state"

type stateJSON struct {
	UUID       string                     `json:"uuid"`
	Benchmarks map[string]*sonm.Benchmark `json:"benchmarks"`
}

func newEmptyState() *stateJSON {
	return &stateJSON{
		UUID:       uuid.New(),
		Benchmarks: map[string]*sonm.Benchmark{},
	}
}

type state struct {
	mu  sync.Mutex
	ctx context.Context
	s   store.Store

	uuid       string
	benchmarks map[string]*sonm.Benchmark
}

func initStorage(p string) (store.Store, error) {
	boltdb.Register()
	config := store.Config{
		Bucket: "sonm",
	}

	return libkv.NewStore(store.BOLTDB, []string{p}, &config)
}

// NewState returns state storage that uses boltdb as backend
func NewState(ctx context.Context, config Config) (*state, error) {
	stor, err := initStorage(config.Store())
	if err != nil {
		return nil, err
	}

	s := &state{
		ctx: ctx,
		s:   stor,

		uuid:       "",
		benchmarks: make(map[string]*sonm.Benchmark),
	}

	err = s.loadInitial()
	if err != nil {
		return nil, err
	}

	return s, err
}

// loadInitial loads state from boltdb
func (s *state) loadInitial() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	kv, err := s.s.Get(stateKey)
	if err != nil && err != store.ErrKeyNotFound {
		return err

	}

	data := &stateJSON{}
	if kv != nil {
		// unmarshal exiting state
		err = json.Unmarshal(kv.Value, &data)
		if err != nil {
			return err
		}
	} else {
		// create new state (clean start)
		data = newEmptyState()
	}

	s.uuid = data.UUID
	s.benchmarks = data.Benchmarks

	err = s.save()
	if err != nil {
		return fmt.Errorf("cannot save state into storage: %v", err)
	}

	return nil
}

// save dumps current state on disk.
//
// Warn: need no be protected by `s.mu` mutex
func (s *state) save() error {
	data := &stateJSON{
		UUID:       s.uuid,
		Benchmarks: s.benchmarks,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return s.s.Put(stateKey, b, &store.WriteOptions{})
}

func (s *state) getID() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.uuid
}

func (s *state) setID(v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.uuid = v
	return s.save()
}

func (s *state) getBenchmarkResults() map[string]*sonm.Benchmark {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.benchmarks
}

func (s *state) setBenchmarkResults(v map[string]*sonm.Benchmark) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.benchmarks = v
	return s.save()
}
