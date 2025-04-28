package main

import (
	"fmt"
	"net/http"
	"os"

	controller "github.com/WoodProgrammer/postgresql-connection-manager/controller"
	lib "github.com/WoodProgrammer/postgresql-connection-manager/lib"

	"github.com/gin-gonic/gin"
)

const (
	CreateCgroupsPath    = "/v1/create-cgroups"
	DeleteCgroupsPath    = "/v1/delete-cgroups"
	Metrics              = "/v1/metrics"
	MovePIDToCgroupsPath = "/v1/move-pid-to-cgroups"
	GetPIDOfQueries      = "/v1/get-pid-of-queries"
)

func NewCgroupHandlerClient() lib.CgroupInterface {
	return &lib.CgroupHandler{}
}

func NewMetricClient() lib.MetricInterface {
	return &lib.MetricHandler{}
}

func NewControllerClient() *controller.Controller {
	c := NewCgroupHandlerClient()
	m := NewMetricClient()
	return &controller.Controller{
		CGroupClient: c,
		MetricClient: m,
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := os.Getenv("PG_CONNECTION_AUTH_TOKEN")
		if len(authToken) == 0 {
			panic("There is no auth token on environment variable, exiting immediately ....") // :D
		}

		tokenPath := fmt.Sprintf("Bearer %s", authToken)
		token := c.GetHeader("Authorization")

		if token == "" || token != tokenPath {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}

func main() {
	port := os.Getenv("PG_CONNECTION_HANDLER_PORT")

	if len(port) == 0 {
		port = "8080"
	}
	router := gin.Default()
	controllerHandler := NewControllerClient()
	router.POST(CreateCgroupsPath, AuthMiddleware(), controllerHandler.CreateCgroup)
	router.POST(MovePIDToCgroupsPath, AuthMiddleware(), controllerHandler.MovePIDToCgroup)
	router.DELETE(DeleteCgroupsPath, AuthMiddleware(), controllerHandler.DeleteCgroupsPath)
	router.GET(GetPIDOfQueries, AuthMiddleware(), controllerHandler.GetPIDOfQueries)
	router.GET(Metrics, AuthMiddleware(), controllerHandler.GetMetrics)

	router.Run(fmt.Sprintf("localhost:%s", port))
}
