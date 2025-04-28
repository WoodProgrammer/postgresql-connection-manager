package controller

import lib "github.com/WoodProgrammer/postgresql-connection-manager/lib"

type Controller struct {
	CGroupClient lib.CgroupInterface
	MetricClient lib.MetricInterface
}
type Metrics struct {
	CgroupName        string
	MemoryStatMetrics MemoryStat
}

type MemoryStat struct {
	Vmalloc string
}
type CGroupV2DeletionRequest struct {
	Name string `json:"name"`
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

type GetPIDOfQueriesRequest struct {
	Query    string `json:"query"`
	Port     string `json:"port"`
	Password string `json:"password"`
	UserName string `json:"username"`
	SSLMode  string `json:"sslmode"`
}
