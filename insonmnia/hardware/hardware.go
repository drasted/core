package hardware

import (
	"fmt"

	"github.com/cnf/structhash"
	"github.com/shirou/gopsutil/mem"
	"github.com/sonm-io/core/insonmnia/hardware/cpu"
	"github.com/sonm-io/core/proto"
)

type CPUProperties struct {
	Device    cpu.Device                 `json:"device"`
	Benchmark map[string]*sonm.Benchmark `json:"benchmark"`
}

type MemoryProperties struct {
	Device    *mem.VirtualMemoryStat     `json:"device"`
	Benchmark map[string]*sonm.Benchmark `json:"benchmark"`
}

type GPUProperties struct {
	Device    *sonm.GPUDevice            `json:"device"`
	Benchmark map[string]*sonm.Benchmark `json:"benchmark"`
}

type NetworkProperties struct {
	Device    interface{}                `json:"device"`
	Benchmark map[string]*sonm.Benchmark `json:"benchmark"`
}

type StorageProperties struct {
	Device    interface{}                `json:"device"`
	Benchmark map[string]*sonm.Benchmark `json:"benchmark"`
}

// Hardware accumulates the finest hardware information about system the worker
// is running on.
type Hardware struct {
	CPU     []*CPUProperties   `json:"cpu"`
	GPU     []*GPUProperties   `json:"gpu"`
	Memory  *MemoryProperties  `json:"memory"`
	Network *NetworkProperties `json:"network"`
	Storage *StorageProperties `json:"storage"`
}

// NewHardware returns initial hardware capabilities for Worker's host.
// Parts of the struct may be filled later by HW-plugins.
func NewHardware() (*Hardware, error) {
	hw := &Hardware{
		CPU:     []*CPUProperties{},
		GPU:     []*GPUProperties{},
		Memory:  &MemoryProperties{Benchmark: make(map[string]*sonm.Benchmark)},
		Network: &NetworkProperties{Benchmark: make(map[string]*sonm.Benchmark)},
		Storage: &StorageProperties{Benchmark: make(map[string]*sonm.Benchmark)},
	}

	CPUs, err := cpu.GetCPUDevices()
	if err != nil {
		return nil, err
	}

	for _, dev := range CPUs {
		hw.CPU = append(hw.CPU, &CPUProperties{
			Device:    dev,
			Benchmark: make(map[string]*sonm.Benchmark),
		})
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	hw.Memory = &MemoryProperties{
		Device:    vm,
		Benchmark: make(map[string]*sonm.Benchmark),
	}

	return hw, nil
}

// LogicalCPUCount returns the number of logical CPUs in the system.
func (h *Hardware) LogicalCPUCount() int {
	count := 0
	for _, c := range h.CPU {
		count += int(c.Device.Cores)
	}

	return count
}

func (h *Hardware) Hash() string {
	return h.devicesMap().Hash()
}

type HashableMemory struct {
	Total uint64 `json:"total"`
}

// DeviceMapping maps hardware capabilities to device description, hashing-friendly
type DeviceMapping struct {
	CPU     []cpu.Device      `json:"cpu"`
	GPU     []*sonm.GPUDevice `json:"gpu"`
	Memory  HashableMemory    `json:"memory"`
	Network interface{}       `json:"network"`
	Storage interface{}       `json:"storage"`
}

func (dm *DeviceMapping) Hash() string {
	return fmt.Sprintf("%x", structhash.Md5(dm, 1))
}

// BenchmarkMapping maps hardware units to benchmark names
// that where passed for this device
//type BenchmarkMapping struct {
//	CPU     []string `json:"cpu"`
//	GPU     []string `json:"gpu"`
//	Memory  []string `json:"memory"`
//	Network []string `json:"network"`
//	Storage []string `json:"storage"`
//}
//
//func (bm *BenchmarkMapping) Hash() string {
//	return fmt.Sprintf("%x", structhash.Md5(bm, 1))
//}

func (h *Hardware) devicesMap() *DeviceMapping {
	m := &DeviceMapping{
		CPU:     []cpu.Device{},
		GPU:     []*sonm.GPUDevice{},
		Memory:  HashableMemory{Total: h.Memory.Device.Total},
		Network: h.Network.Device,
		Storage: h.Storage.Device,
	}

	for _, c := range h.CPU {
		m.CPU = append(m.CPU, c.Device)
	}

	for _, g := range h.GPU {
		m.GPU = append(m.GPU, g.Device)
	}

	return m
}

//func (h *Hardware) benchmarkMap() *BenchmarkMapping {
//	m := &BenchmarkMapping{
//		CPU:     []string{},
//		GPU:     []string{},
//		Memory:  []string{},
//		Network: []string{},
//		Storage: []string{},
//	}
//
//	for _, c := range h.CPU {
//		for id := range c.Benchmark {
//			m.CPU = append(m.CPU, id)
//		}
//	}
//
//	for _, g := range h.GPU {
//		for id := range g.Benchmark {
//			m.GPU = append(m.GPU, id)
//		}
//	}
//
//	for id := range h.Memory.Benchmark {
//		m.Memory = append(m.Memory, id)
//	}
//
//	for id := range h.Network.Benchmark {
//		m.Network = append(m.Network, id)
//	}
//
//	for id := range h.Storage.Benchmark {
//		m.Storage = append(m.Storage, id)
//	}
//
//	return m
//}

//func (h *Hardware) benchmarkTypes() map[string]bool {
//	m := map[string]bool{}
//	for _, c := range h.CPU {
//		for id := range c.Benchmark {
//			m[id] = true
//		}
//	}
//
//	for _, g := range h.GPU {
//		for id := range g.Benchmark {
//			m[id] = true
//		}
//	}
//
//	for id := range h.Memory.Benchmark {
//		m[id] = true
//	}
//
//	for id := range h.Network.Benchmark {
//		m[id] = true
//	}
//
//	for id := range h.Storage.Benchmark {
//		m[id] = true
//	}
//
//	return m
//}
