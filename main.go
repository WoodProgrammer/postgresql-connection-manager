package main

import (
	controller "github.com/WoodProgrammer/postgresql-connection-manager/controller"
	lib "github.com/WoodProgrammer/postgresql-connection-manager/lib"

	"github.com/gin-gonic/gin"
)

const (
	CreateCgroupsPath    = "/v1/create-cgroups"
	MovePIDToCgroupsPath = "/v1/move-pid-to-cgroups"
	GetPIDOfQueries      = "/v1/get-pid-of-queries"
)

func NewCgroupHandlerClient() lib.CgroupInterface {
	return &lib.CgroupHandler{}
}

func NewControllerClient() *controller.Controller {
	c := NewCgroupHandlerClient()
	return &controller.Controller{
		CGroupClient: c,
	}
}

func main() {

	router := gin.Default()
	controllerHandler := NewControllerClient()
	router.POST(CreateCgroupsPath, controllerHandler.CreateCgroup)
	router.POST(MovePIDToCgroupsPath, controllerHandler.MovePIDToCgroup)
	router.GET(GetPIDOfQueries, controllerHandler.GetPIDOfQueries)

	router.Run("localhost:8080")
}
