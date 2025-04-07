English | [简体中文](./README-cn.md)

# Cache Applications
This page can be read in conjunction with [Caching Applications](https://en.dtm.pub/app/cache.html) in the dtm documentation

## Overview
This project demonstrates how [dtm-labs/rockscache](https://github.com/dtm-labs/rockscache) and [dtm-labs/dtm](https://github.com/dtm-labs/dtm) together can be used to maintain cache consistency, including the following
- Ensuring eventual consistency: demonstrates how rockscache can ensure the consistency between cache and DB
- Ensuring atomic operations: if a process crashes after updating the DB, the cache can still be updated
- Other features of cache management: demonstrates a delayed-delete cache management method that is anti-penetration, anti-breakdown and anti-avalanche
- Strong consistency usage: demonstrates how to provide strong consistency to application
- Strong Consistency in Upgrades and Downgrades: provide strong consistency to application even if the cache is being downgraded and upgraded

## Start dtm
[Quick start dtm](https://en.dtm.pub/guide/install.html)

## Run This Example
`go run main.go`

## Ensure Eventual Consistency
The code is mainly in demo/api-version, and the example demonstrates the combination of rockscache and dtm to solve the version inconsistency problem that the delete cache fails to.
- A request is made to demonstrate the version inconsistency problem with the normal delete cache solution `curl http://localhost:8081/api/busi/version?mode=delete`
- Launching a request to demonstrate rockscache+dtm to ensure version inconsistency `curl http://localhost:8081/api/busi/version?mode=rockscache`

In this demo, the main points are as follows
1. It will initialize the DB data to v1, then query the cache, query the DB to get v1 after no data, sleep for a few seconds to simulate network latency, and then update the cache.
2. Will then modify the data to v2, then query the cache, query the DB for v2 after no data, sleep for a few milliseconds to simulate no network latency, and then update the cache.

- For mode=delete, the network latency causes v1 to be written to the cache after v2, resulting in the final version in the cache being v1, which is not consistent with the database
- For mode=rockscache, although network latency is encountered, the final version in the cache is v2, which is the same as in the database
### Guaranteed atomicity
- Initiate a normal request to update data, simulating a crash that causes the DB to be inconsistent with the cache `curl http://localhost:8081/api/busi/atomic?mode=none`
- Initiate a request to update data via dtm, simulating a crash, resulting in the DB being inconsistent with the cache, but after 5s the cache is updated to the same value as DB `curl http://localhost:8081/api/busi/atomic?mode=dtm`

You can see that updating the cache via dtm can ensure that both the DB and cache update succeed or fail in the event of a process crash.

This example also demonstrates that dtm can be used not only in combination with rockscache, but also independently of rockscache. The architectural solution is simpler than local message tables, transaction messages, binlog listeners, etc.

### Strong consistent access
rockscache+dtm can also provide strongly consistent access, in the same way as the client of rockscache, by setting the `StringConsistency` option

Launch an example demonstrating strong consistent access `curl http://localhost:8081/api/busi/strong`, which does the following
1. initialise the data, including the database and cache
2. update the data using a dtm 2-phase message to ensure atomicity of the DB and cache operations, where the cached data takes 3s to compute
3. query the global transaction status before the cache update is completed, the status is not completed. (If the user queries the business result, even if the DB update is completed, as long as the global transaction is not completed, the user should be informed that it is not completed)
4. query cache, in strong consistent access mode, this query will wait for the result in 2, different from the eventual consistency


### Strong Consistency on Downgrading and Upgrading
If strong consistency is also required for the short time window of downgrading and upgrading, rockscache+dtm can do also.

Launch an example `curl http://localhost:8081/api/busi/downgrade` that demonstrates strong consistency even during downgrading and upgrading, which does the following:
1. initialize the data, including the database and cache
2. use dtm's Saga mode to update the data, ensuring atomicity for the three operations of locking the cache, updating the DB, and deleting the cache, in addition to the 3s computation time required for the cache data
3. after updating the DB, before updating the cache, you can tell the user that the operation has been completed (different from the previous case of strong consistency, the condition is relaxed)
4. query cache, in strong consistent access mode, the query will wait for the result in 2
5. when downgrading, turn off the read cache first, wait until all reads are not accessing the cache, then turn off the delete cache
6. When upgrading, first turn on the delete cache to ensure that all DB updates are updated to the cache, and then turn on the read cache

Note: In 2, we use SAGA transaction mode instead of message because we have to consider that update DB may fail.
