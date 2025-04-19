package lib

import (
	"fmt"
	"os"

	cgroupsv2 "github.com/containerd/cgroups/v3/cgroup2"
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
		fmt.Printf("Error creating cgroup: %v\n", err)
		return err
	} else {
		fmt.Println("The group created successfully")
	}
	fmt.Println(cgroupManager)

	return nil
}

func (c *CgroupHandler) MovePIDToCgroupHandler(name string, pid string) error {
	cgroupPath := "/sys/fs/cgroup/"
	content := []byte(pid)

	err := os.WriteFile(cgroupPath+"/"+name, content, 0644)
	if err != nil {
		panic(err)
	}
	return nil

}
