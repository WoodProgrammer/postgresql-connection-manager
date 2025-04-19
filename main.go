package main

import (
	controller "github.com/WoodProgrammer/postgresql-connection-manager/controller"

	"github.com/gin-gonic/gin"
)

const (
	CreateCgroupsPath    = "/v1/create-cgroups"
	MovePIDToCgroupsPath = "/v1/move-pid-to-cgroups"
)

func main() {
	router := gin.Default()

	router.POST(CreateCgroupsPath, controller.CreateCgroup)
	router.POST(MovePIDToCgroupsPath, controller.MovePIDToCgroup)

	router.Run("localhost:8080")
}
