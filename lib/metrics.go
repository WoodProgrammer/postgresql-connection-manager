package lib

import (
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

type CgroupCollector struct {
	Desc *prometheus.Desc
}

func (c *CgroupCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Desc
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

func (m *CgroupCollector) Collect(ch chan<- prometheus.Metric) {

	directoryList, err := os.ReadDir(CGROUP_PATH)
	if err != nil {
		log.Err(err).Msgf("Error while inspecting in cgroup path")
	}

	for _, dir := range directoryList {
		if strings.Contains(dir.Name(), "pg-") {
			dirName := dir.Name()
			cpuUsage, cpuThrottled, err := parseCPUStat(CGROUP_PATH + "/" + dirName + "/cpu.stat")
			if err != nil {
				log.Err(err).Msgf("There is and error while reading cpu stats")
			}

			ch <- prometheus.MustNewConstMetric(
				m.Desc,
				prometheus.GaugeValue,
				float64(cpuUsage),
				"cpu_usage", dirName,
			)

			ch <- prometheus.MustNewConstMetric(
				m.Desc,
				prometheus.GaugeValue,
				float64(cpuThrottled),
				"cpu_throttled", dirName,
			)

			memCurrent, err := readUintFromFile(CGROUP_PATH + "/" + dirName + "/memory.current")
			if err != nil {
				log.Err(err).Msgf("There is and error while reading memory.current")
			}
			ch <- prometheus.MustNewConstMetric(
				m.Desc,
				prometheus.GaugeValue,
				float64(memCurrent),
				"mem_current", dirName,
			)

			oomEvents, err := parseCgroupEvents(CGROUP_PATH + "/" + dirName + "/cgroup.events")
			if err != nil {
				log.Err(err).Msgf("There is and error while reading cgroup.events")

			}
			ch <- prometheus.MustNewConstMetric(
				m.Desc,
				prometheus.GaugeValue,
				float64(oomEvents),
				"oom_events", dirName,
			)

			pidsCurrent, err := readUintFromFile(CGROUP_PATH + "/" + dirName + "/pids.current")
			if err != nil {
				log.Err(err).Msgf("There is and error while reading pids.current")
			}

			ch <- prometheus.MustNewConstMetric(
				m.Desc,
				prometheus.GaugeValue,
				float64(pidsCurrent),
				"pids.current", dirName,
			)
		} else {
			continue
		}
	}

}
