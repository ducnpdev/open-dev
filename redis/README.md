# Redis
## Installation Demo Simple 
- read article: https://viblo.asia/p/golang-redis-luu-don-gian-Az45bjNO5xY

## build redis cluster


## use helm 
ducnp@nguyens-MacBook-Pro-4 redis % helm install redis-cache bitnami/redis -f custom-values.yaml 
NAME: redis-cache
LAST DEPLOYED: Sat Nov 18 11:43:36 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: redis
CHART VERSION: 18.3.3
APP VERSION: 7.2.3

** Please be patient while the chart is being deployed **

Redis&reg; can be accessed via port 6379 on the following DNS name from within your cluster:

    redis-cache-master.default.svc.cluster.local



To get your password run:

    export REDIS_PASSWORD=$(kubectl get secret --namespace default redis-cache -o jsonpath="{.data.redis-password}" | base64 -d)

To connect to your Redis&reg; server:

1. Run a Redis&reg; pod that you can use as a client:

   kubectl run --namespace default redis-client --restart='Never'  --env REDIS_PASSWORD=$REDIS_PASSWORD  --image docker.io/bitnami/redis:7.2.3-debian-11-r1 --command -- sleep infinity

   Use the following command to attach to the pod:

   kubectl exec --tty -i redis-client \
   --namespace default -- bash

2. Connect using the Redis&reg; CLI:
   REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h redis-cache-master

To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace default svc/redis-cache-master 6379:6379 &
    REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h 127.0.0.1 -p 6379

### in k8s
 kubectl exec -it redis-cache-master-0 -- sh 
  REDISCLI_AUTH=g4m3C4Rizt redis-cli
redis-cache-master-0

# Memory
 <!-- INFO MEMORY -->
used_memory:75726648
used_memory_human:72.22M
used_memory_rss:84209664
used_memory_rss_human:80.31M
used_memory_peak:76287104
used_memory_peak_human:72.75M
used_memory_peak_perc:99.27%
used_memory_overhead:1416168
used_memory_startup:898976
used_memory_dataset:74310480
used_memory_dataset_perc:99.31%
allocator_allocated:76079992
allocator_active:76333056
allocator_resident:82345984
total_system_memory:8233017344
total_system_memory_human:7.67G
used_memory_lua:31744
used_memory_vm_eval:31744
used_memory_lua_human:31.00K
used_memory_scripts_eval:0
number_of_cached_scripts:0
number_of_functions:0
number_of_libraries:0
used_memory_vm_functions:32768
used_memory_vm_total:64512
used_memory_vm_total_human:63.00K
used_memory_functions:184
used_memory_scripts:184
used_memory_scripts_human:184B
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
allocator_frag_ratio:1.00
allocator_frag_bytes:253064
allocator_rss_ratio:1.08
allocator_rss_bytes:6012928
rss_overhead_ratio:1.02
rss_overhead_bytes:1863680
mem_fragmentation_ratio:1.11
mem_fragmentation_bytes:8503688
mem_not_counted_for_evict:8
mem_replication_backlog:0
mem_total_replication_buffers:0
mem_clients_slaves:0
mem_clients_normal:1928
mem_cluster_links:0
mem_aof_buffer:8
mem_allocator:jemalloc-5.3.0
active_defrag_running:0
lazyfree_pending_objects:0
lazyfreed_objects:0


volatile-lru - remove a key among the ones with an expire set, trying to remove keys not recently used.
volatile-ttl - remove a key among the ones with an expire set, trying to remove keys with short remaining time to live.
volatile-random - remove a random key among the ones with an expire set.
allkeys-lru - like volatile-lru, but will remove every kind of key, both normal keys or keys with an expire set.
allkeys-random - like volatile-random, but will remove every kind of keys, both normal keys and keys with an expire set.
--->
volatile-lru - Loại bỏ một khóa trong số những khóa có đặt thời gian hết hạn, cố gắng loại bỏ những khóa không được sử dụng gần đây.
volatile-ttl - Loại bỏ một khóa trong số những khóa có đặt thời gian hết hạn, cố gắng loại bỏ những khóa có thời gian sống còn lại ngắn.
volatile-random - Loại bỏ một khóa ngẫu nhiên trong số những khóa có đặt thời gian hết hạn.
allkeys-lru - Giống như volatile-lru, nhưng sẽ loại bỏ mọi loại khóa, cả khóa thông thường và khóa có đặt thời gian hết hạn.
allkeys-random - Giống như volatile-random, nhưng sẽ loại bỏ mọi loại khóa, cả khóa thông thường và khóa có đặt thời gian hết hạn.

## Pub sub
- code example pub-sub redis