package controller

import (
	lib "github.com/WoodProgrammer/postgresql-connection-manager/lib"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	CGroupClient lib.CgroupInterface
}

type CGroupV2CreationRequest struct {
	Name      string `json:"name"`
	PID       string `json:"pid"`
	CpuCycle  int64  `json:"cycle"`
	CpuPeriod int64  `json:"period"`
	Memory    int64  `json:"memory"`
}

type CGroupV2MoveRequest struct {
	PID  string `json:"pid"`
	Name string `json:"name"`
}

func (c *Controller) CreateCgroup(ctx *gin.Context) {
	var cgroup CGroupV2CreationRequest

	if err := ctx.BindJSON(&cgroup); err != nil {
		return
	}
	res := c.CGroupClient.HandleCgroupResources(cgroup.CpuCycle, cgroup.Memory, uint64(cgroup.CpuPeriod))
	err := c.CGroupClient.CreateCgroupV2(res, "", cgroup.Name)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.JSON(200, &cgroup)
}

func (c *Controller) MovePIDToCgroup(ctx *gin.Context) {
	var cgroup CGroupV2MoveRequest

	if err := ctx.BindJSON(&cgroup); err != nil {
		return
	}
	err := c.CGroupClient.MovePIDToCgroupHandler(cgroup.Name, cgroup.PID)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.JSON(200, &cgroup)
}
