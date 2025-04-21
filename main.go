package main

import (
	"fmt"
	"os"

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
	port := os.Getenv("PG_CONNECTION_HANDLER_PORT")

	if len(port) == 0 {
		port = "8080"
	}
	router := gin.Default()
	controllerHandler := NewControllerClient()
	router.POST(CreateCgroupsPath, controllerHandler.CreateCgroup)
	router.POST(MovePIDToCgroupsPath, controllerHandler.MovePIDToCgroup)
	router.GET(GetPIDOfQueries, controllerHandler.GetPIDOfQueries)

	router.Run(fmt.Sprintf("localhost:%s", port))
}
