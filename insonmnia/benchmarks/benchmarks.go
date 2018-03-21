package benchmarks

import "github.com/sonm-io/core/proto"

type BenchList interface {
	List() (map[string]*sonm.Benchmark, error)
}

type dumbBenchmark struct{}

// TODO(sshaman1101): make lockable wathcer for bench list (like Hub's whitelist do)

func NewDumbBenchmarks() BenchList {
	return &dumbBenchmark{}
}

func (db *dumbBenchmark) List() (map[string]*sonm.Benchmark, error) {
	return map[string]*sonm.Benchmark{
		"test-cpu": {ID: "test-cpu", Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_CPU},
		"test-ram": {ID: "test-ram", Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_RAM},
		"test-gpu": {ID: "test-gpu", Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_GPU},
	}, nil
}
