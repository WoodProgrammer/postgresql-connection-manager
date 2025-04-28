package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (c *Controller) GetMetrics(ctx *gin.Context) {

	cgroupName := ctx.Query("cgroupName")
	metrics, err := c.MetricClient.CollectMetrics(cgroupName)

	if err != nil {
		log.Err(err).Msgf("Error while fetching metrics controller.GetMetrics()")
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, &metrics)
}
func (c *Controller) DeleteCgroupsPath(ctx *gin.Context) {
	var req CGroupV2DeletionRequest

	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	err := c.CGroupClient.DeleteGroupV2(req.Name)
	if err != nil {
		log.Err(err).Msgf("Error while deleting cgroups controller.CreateCgroup()")
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, "OK")
}

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

func (c *Controller) GetPIDOfQueries(ctx *gin.Context) {
	var request GetPIDOfQueriesRequest

	if err := ctx.BindJSON(&request); err != nil {
		return
	}

	result, err := c.CGroupClient.GatherPostgresqlConnectionDetails("localhost", request.Port, request.Password, request.UserName, request.SSLMode, request.Query)
	if err != nil {
		log.Err(err).Msgf("Error while fetching query results in postgresql")
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, result)
}
