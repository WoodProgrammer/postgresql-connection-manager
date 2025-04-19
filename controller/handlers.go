package controller

import (
	lib "github.com/WoodProgrammer/postgresql-connection-manager/lib"
	"github.com/gin-gonic/gin"
)

func NewCgroupHandlerClient() lib.CgroupInterface {
	return &lib.CgroupHandler{}
}

type CGroupV2CreationRequest struct {
	Name      string `json:"name"`
	PID       string `json:"pid"`
	CpuCycle  int64  `json:"cycle"`
	CpuPeriod int64  `json:"period"`
	Memmory   int64  `json:"memory"`
}

type CGroupV2MoveRequest struct {
	PID  string `json:"pid"`
	Name string `json:"name"`
}

func CreateCgroup(c *gin.Context) {
	cgroupHandlerClient := NewCgroupHandlerClient()
	var cgroup CGroupV2CreationRequest

	if err := c.BindJSON(&cgroup); err != nil {
		return
	}
	res := cgroupHandlerClient.HandleCgroupResources(cgroup.CpuCycle, cgroup.CpuCycle, uint64(cgroup.Memmory))
	err := cgroupHandlerClient.CreateCgroupV2(res, "", cgroup.Name)
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, &cgroup)
}

func MovePIDToCgroup(c *gin.Context) {
	var cgroup CGroupV2MoveRequest

	if err := c.BindJSON(&cgroup); err != nil {
		return
	}
	cgroupHandlerClient := NewCgroupHandlerClient()
	err := cgroupHandlerClient.MovePIDToCgroupHandler(cgroup.Name, cgroup.PID)
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, &cgroup)
}
