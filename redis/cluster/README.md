# redis cluster

docker run -d -it -v `pwd`/7000:/redis --name node-1 -p 6380:7000 --network redis-cluster redis redis-server /redis/redis.conf
docker run -d -it -v `pwd`/7001:/redis --name node-2 -p 6381:7000 --network redis-cluster redis redis-server /redis/redis.conf
docker run -d -it -v `pwd`/7002:/redis --name node-3 -p 6382:7000 --network redis-cluster redis redis-server /redis/redis.conf