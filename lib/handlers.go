package lib

import (
	"os"

	cgroupsv2 "github.com/containerd/cgroups/v3/cgroup2"
	"github.com/rs/zerolog/log"
)

type CgroupInterface interface {
	HandleCgroupResources(cpuQuota, memoryInByes int64, period uint64) cgroupsv2.Resources
	CreateCgroupV2(res cgroupsv2.Resources, cgroupPath string, cgroupName string) error
	MovePIDToCgroupHandler(name string, pid string) error
}

type CgroupHandler struct {
}

func (c *CgroupHandler) HandleCgroupResources(cpuQuota, memoryInByes int64, period uint64) cgroupsv2.Resources {
	res := cgroupsv2.Resources{}

	max := cgroupsv2.NewCPUMax(&cpuQuota, &period)
	cpu := cgroupsv2.CPU{Max: max}
	memory := cgroupsv2.Memory{Max: &memoryInByes}
	res = cgroupsv2.Resources{CPU: &cpu, Memory: &memory}
	return res

}

func (c *CgroupHandler) CreateCgroupV2(res cgroupsv2.Resources, cgroupPath string, cgroupName string) error {
	if len(cgroupPath) == 0 {
		cgroupPath = "/sys/fs/cgroup/"
	}
	cgroupManager, err := cgroupsv2.NewManager(cgroupPath, "/"+cgroupName, &res)
	if err != nil {
		log.Err(err).Msg("Error creating cgroup: in cGroupHandler CreateCgroupV2")
		return err
	}
	log.Info().Msgf("The group created successfully %s object is %s", cgroupName, cgroupManager)
	return nil
}

func (c *CgroupHandler) MovePIDToCgroupHandler(name string, pid string) error {
	cgroupPath := "/sys/fs/cgroup/"
	content := []byte(pid)

	err := os.WriteFile(cgroupPath+name+"/cgroup.procs", content, 0644)
	if err != nil {
		log.Err(err).Msg("Error while assign the PID to the relevant cgroups MovePIDToCgroupHandler")
		return err
	}
	return nil

}
