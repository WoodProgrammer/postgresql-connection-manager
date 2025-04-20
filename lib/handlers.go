package lib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	cgroupsv2 "github.com/containerd/cgroups/v3/cgroup2"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type CgroupInterface interface {
	HandleCgroupResources(cpuQuota, memoryInByes int64, period uint64) cgroupsv2.Resources
	CreateCgroupV2(res cgroupsv2.Resources, cgroupName string) error
	MovePIDToCgroupHandler(name string, pid string) error
	GatherPostgresqlConnectionDetails(host, port, password, user, sslmode, query string) (string, error)
}

type CgroupHandler struct {
}

func (c *CgroupHandler) HandleCgroupResources(cpuQuota, memoryInByes int64, period uint64) cgroupsv2.Resources {
	res := cgroupsv2.Resources{}

	max := cgroupsv2.NewCPUMax(&cpuQuota, &period)
	cpu := cgroupsv2.CPU{Max: max}
	memory := cgroupsv2.Memory{Max: &memoryInByes}
	res = cgroupsv2.Resources{CPU: &cpu, Memory: &memory}
	return res

}

func (c *CgroupHandler) CreateCgroupV2(res cgroupsv2.Resources, cgroupName string) error {

	cgroupManager, err := cgroupsv2.NewManager(CGROUP_PATH, "/"+cgroupName, &res)
	if err != nil {
		log.Err(err).Msg("Error creating cgroup: in cGroupHandler CreateCgroupV2")
		return err
	}
	log.Info().Msgf("The group created successfully %s object is %s", cgroupName, cgroupManager)
	return nil
}

func (c *CgroupHandler) MovePIDToCgroupHandler(name string, pid string) error {
	content := []byte(pid)

	err := os.WriteFile(CGROUP_PATH+name+"/cgroup.procs", content, FILE_PERMISSION)
	if err != nil {
		log.Err(err).Msg("Error while assign the PID to the relevant cgroups MovePIDToCgroupHandler")
		return err
	}
	return nil
}

func (c *CgroupHandler) GatherPostgresqlConnectionDetails(host, port, password, user, sslmode, query string) (string, error) {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s", host, port, user, password, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Err(err).Msg("Error while open connection to postgresql database")
		return "", err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Err(err).Msg("Cannot connect to DB:")
		return "", err
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Err(err).Msg("Error while executing query")
		return "", err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Err(err).Msg("Columns error:")
		return "", err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Allocate slice for row values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the row
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Err(err).Msg("Scan error:")
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}
		results = append(results, rowMap)
	}

	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Err(err).Msg("JSON marshal error:")
		return "", err
	}

	return string(jsonBytes), nil

}
