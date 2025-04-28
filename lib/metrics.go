package lib

import (
	"os"
	"strconv"
	"strings"
)

type MetricInterface interface {
	CollectMetrics(cgroupPath string) (*CgroupMetrics, error)
}

type MetricHandler struct {
}

type CgroupMetrics struct {
	CPUUsageMicros     uint64 `json:"cpu_usage_micros"`
	CPUThrottledMicros uint64 `json:"cpu_throttled_micros"`
	MemoryCurrent      uint64 `json:"memory_current_bytes"`
	OOMEvents          uint64 `json:"oom_events"`
	ProcessesCurrent   uint64 `json:"processes_current"`
}

func readUintFromFile(filepath string) (uint64, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

func parseCPUStat(filepath string) (usage, throttled uint64, err error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		switch fields[0] {
		case "usage_usec":
			usage, _ = strconv.ParseUint(fields[1], 10, 64)
		case "throttled_usec":
			throttled, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}
	return usage, throttled, nil
}

func parseCgroupEvents(filepath string) (oomEvents uint64, err error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		if fields[0] == "oom" {
			oomEvents, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}
	return oomEvents, nil
}

func (m *MetricHandler) CollectMetrics(cgroupName string) (*CgroupMetrics, error) {
	fullCgroupPath := CGROUP_PATH + cgroupName
	cpuUsage, cpuThrottled, err := parseCPUStat(fullCgroupPath + "/cpu.stat")
	if err != nil {
		return nil, err
	}

	memCurrent, err := readUintFromFile(fullCgroupPath + "/memory.current")
	if err != nil {
		return nil, err
	}

	oomEvents, err := parseCgroupEvents(fullCgroupPath + "/cgroup.events")
	if err != nil {
		return nil, err
	}

	pidsCurrent, err := readUintFromFile(fullCgroupPath + "/pids.current")
	if err != nil {
		return nil, err
	}

	return &CgroupMetrics{
		CPUUsageMicros:     cpuUsage,
		CPUThrottledMicros: cpuThrottled,
		MemoryCurrent:      memCurrent,
		OOMEvents:          oomEvents,
		ProcessesCurrent:   pidsCurrent,
	}, nil
}
