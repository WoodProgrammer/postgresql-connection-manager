package controller

import lib "github.com/WoodProgrammer/postgresql-connection-manager/lib"

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
