# postgresql-connection-manager
This is project to manage postgresql connections via cgroup V2



curl http://localhost:8080/v1/create-cgroups \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"pid": "4","name": "pg-exporter-cgroup","period":1000, "cycle": 1000, "memory": 1000}'