package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (c *Controller) CreateCgroup(ctx *gin.Context) {
	var cgroup CGroupV2CreationRequest

	if err := ctx.BindJSON(&cgroup); err != nil {
		return
	}
	res := c.CGroupClient.HandleCgroupResources(cgroup.CpuCycle, cgroup.Memory, uint64(cgroup.CpuPeriod))
	err := c.CGroupClient.CreateCgroupV2(res, cgroup.Name)
	if err != nil {
		log.Err(err).Msgf("Error while creating cgroups controller.CreateCgroup()")
		ctx.JSON(500, err)
		return
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
		log.Err(err).Msgf("Error while moving PID to %s controller.MovePIDToCgroup()", cgroup.PID)
		ctx.JSON(500, err)
		return
	}
	ctx.JSON(200, &cgroup)
}
