package benchmarks

import (
	"github.com/sonm-io/core/proto"
)

const (
	// can get from system
	CPUCores    = "cpu-cores"
	RamSize     = "ram-size"
	StorageSize = "storage-size"
	GPUCount    = "gpu-count"
	GPUMem      = "gpu-mem"

	// must be measured by container
	CPUSysbenchMulti  = "cpu-sysbench-multi"
	CPUSysbenchSingle = "cpu-sysbench-single"
	NetDownload       = "net-download"
	NetUpload         = "net-upload"
	GPUEthHashrate    = "gpu-eth-hashrate"
	GPUCashHashrate   = "gpu-cash-hashrate"
	GPURedshift       = "gpu-redshift"
)

var specialBenchmarks = []string{CPUCores, RamSize, StorageSize, GPUCount, GPUMem}

func IsSpecialBenchmark(b *sonm.Benchmark) bool {
	for _, s := range specialBenchmarks {
		if s == b.GetID() {
			return true
		}
	}

	return false
}

type BenchList interface {
	List() map[string]*sonm.Benchmark
}

type dumbBenchmark struct{}

func NewDumbBenchmarks() BenchList {
	return &dumbBenchmark{}
}

func (db *dumbBenchmark) List() map[string]*sonm.Benchmark {
	return map[string]*sonm.Benchmark{
		CPUCores:    {ID: CPUCores, Type: sonm.DeviceType_DEV_CPU},
		RamSize:     {ID: RamSize, Type: sonm.DeviceType_DEV_RAM},
		StorageSize: {ID: StorageSize, Type: sonm.DeviceType_DEV_STORAGE},
		GPUCount:    {ID: GPUCount, Type: sonm.DeviceType_DEV_GPU},
		GPUMem:      {ID: GPUMem, Type: sonm.DeviceType_DEV_GPU},

		CPUSysbenchMulti:  {ID: CPUSysbenchMulti, Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_CPU},
		CPUSysbenchSingle: {ID: CPUSysbenchSingle, Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_CPU},
		NetDownload:       {ID: NetDownload, Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_NETWORK},
		NetUpload:         {ID: NetUpload, Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_NETWORK},
		GPUEthHashrate:    {ID: GPUEthHashrate, Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_GPU},
		GPUCashHashrate:   {ID: GPUCashHashrate, Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_GPU},
		GPURedshift:       {ID: GPURedshift, Image: "sshaman1101/sonm-pidor-bench", Type: sonm.DeviceType_DEV_GPU},
	}
}
