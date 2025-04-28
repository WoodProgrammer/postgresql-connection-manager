# pg-connection-manager

**A tool to isolate PostgreSQL connections using Linux cgroups.** 

This project allows system administrators to detect specific PostgreSQL queries, extract their associated PIDs, and move those processes into dedicated cgroups for fine-grained CPU and memory control.

---

## 🚀 Features

- Detects active PostgreSQL connections based on query filters
- Creates new cgroups (cgroups v2 supported)
- Moves PostgreSQL backend PIDs into specified cgroups
- Enables CPU and memory resource throttling
- Provides an HTTP API for automation and control
- AuthenticationMiddleWare enabled (v0.0.3)
- Prometheus Metrics exposed (v0.0.4)

---

## 🔧 Installation

```bash
git clone https://github.com/WoodProgrammer/pg-connection-manager.git
cd pg-connection-manager
go build -o pg-cgroup-manager

export PG_CONNECTION_HANDLER_PORT=9001
export PG_CONNECTION_AUTH_TOKEN=enc_S1UP3RS3CR3T_4UTH_TOK3n
./pg-connection-manager

```
Please specify the port that you would like to expose.

<hr>

# 📡 API Endpoints

## 1. Get PIDs of Queries
Returns the PIDs of PostgreSQL backend processes that match a given query pattern.

```sh
GET /v1/get-pid-of-queries
```

Basically you can gather PIDs by this endpoint and please specify your query to find out;

### Payload

```sh
curl http://localhost:8080/v1/get-pid-of-queries \
--include \
--header "Authorization: Bearer enc_S1UP3RS3CR3T_4UTH_TOK3n" \
--request "GET" \
--data '{"query": "SELECT pid, usename, application_name, state FROM pg_stat_activity;","port": "5432", "password":"CVVVVV", "username": "postgres", "sslmode": "disable"}'
```

### Sample Response;

```json
[{
    "application_name":"",
        "pid":12416,
        "state":null,
        "usename":null
    },
    ...
```

Then you can create a cgroupV2 and move your postgresql process under the CgroupV2 in specified resource limitations.


<hr></hr>

## 2. Create CgroupsV2

Creates a new cgroup under the default cgroups v2 hierarchy.

```sh
POST /v1/create-cgroups
```

### Payload

There are only three allowed resource groups for now;

* Cpu cycles;
* Cpu period per cycle;
* Memory.max 


To calcuate better values please check the documentation.

```sh
curl http://localhost:8080/v1/create-cgroups \
    --include \
    --header "Authorization: Bearer enc_S1UP3RS3CR3T_4UTH_TOK3n" \
    --request "POST" \
    --data '{"name": "pg-new-cgroup","period":1000, "cycle": 1000, "memory": 536870912}'
```

It basically create cgroup please check the /sys/fs/cgroup directory then you receive 200 OK response


## 3. Move PID to Cgroup
Moves the picked PIDs as you can see on previous payloads (payload-1) into the given cgroup.


```sh
POST /v1/move-pid-to-cgroups
```

### Payload

```sh
curl http://localhost:8080/v1/move-pid-to-cgroups \
    --include \
    --header "Authorization: Bearer enc_S1UP3RS3CR3T_4UTH_TOK3n" \
    --request "POST" \
    --data '{"pid": "7323","name": "pg-long-running-group"}'
```

Then you can basically check the cgroups.procs file of the given cgroup then your postgresql process will runs in given Cgroup.


<hr>

### Metrics

For metrics you have to adjust your prometheus configs like this;

```yaml
scrape_configs:
  - job_name: "pg_cgroup_manager"
    metrics_path: "/v1/metrics"   # change if your exporter exposes metrics on a different path
    scheme: "http"            # or "http" depending on your exporter
    static_configs:
      - targets:
        - "localhost:8080"  # your exporter IP:port or domain

    authorization:
      type: Bearer
      credentials: enc_S1UP3RS3CR3T_4UTH_TOK3n
```

## Dashboard

This project also have very nice dashboard to show up Postgresql connections by groupsV2;

<img src="./img/dashboard.png"></img>

## 🧪 Use Cases
Isolate and throttle heavy or suspicious queries

Enforce resource limits on multi-tenant PostgreSQL instances

Perform controlled performance experiments under constrained resources

Integrate with observability tools for resource-aware database tuning (FUTURE)

<hr>

## 🔐 Security Notes
This tool requires root privileges to interact with the cgroup subsystem.

Ensure the API is protected in production environments (e.g., behind a firewall or with token-based authentication).

<hr>

## 📌 Future Improvements
* Add support for deleting or listing existing cgroups

<hr>

## 🧡 Contributions

Contributions, ideas, and improvements are welcome. If you’re interested in making this tool better, feel free to open an issue or a pull request!
