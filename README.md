# postgresql-connection-manager

This is project to manage postgresql connections via cgroup V2

curl http://localhost:8080/v1/create-cgroups \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "pg-exporter-cgroup-second","period":1000, "cycle": 1000, "memory": 536870912}'


detect pids of queries

SELECT pid, usename, application_name, state FROM pg_stat_activity;

SELECT * FROM generate_series(1, 100) a(id) JOIN generate_series(1, 100) b(id) ON a.id = b.id;


curl http://localhost:8080/v1/move-pid-to-cgroups \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"pid": "12416","name": "pg-exporter-cgroup-second"}'


[45792.694947] oom-kill:constraint=CONSTRAINT_MEMCG,nodemask=(null),cpuset=/,mems_allowed=0,oom_memcg=/pg-exporter-cgroup,task_memcg=/pg-exporter-cgroup,task=postgres,pid=11700,uid=109
[45792.694961] Memory cgroup out of memory: Killed process 11700 (postgres) total-vm:224872kB, anon-rss:3324kB, file-rss:9176kB, shmem-rss:3472kB, UID:109 pgtables:164kB oom_score_adj:0