简体中文 | [English](./README.md)

# 缓存应用
此项目可以结合dtm文档中的 [缓存应用](https://dtm.pub/app/cache.html)阅读

## 概述
本项目主要演示了[dtm-labs/rockscache](https://github.com/dtm-labs/rockscache)和[dtm-labs/dtm](https://github.com/dtm-labs/dtm)配合如何用于维护缓存一致性，包括以下几个方面
- 保证最终一致性：演示了rockscache如何彻底解决缓存与DB版本一致的问题
- 保证原子操作：如果更新完DB后，发生进程crash，也能够保证缓存能够更新
- 缓存管理的其他功能：演示了延迟删除的缓存管理方法，能够防击穿，防穿透，防雪崩（未演示）
- 强一致用法：演示了如何保证强一致
- 升降级中的强一致：如果在缓存的降级和升级中，保证应用访问数据是强一致的

## 启动dtm
[快速启动dtm](https://dtm.pub/guide/install.html)

## 运行本例子
`go run main.go`

### 保证最终一致性
代码主要在demo/api-version，例子主要演示了rockscache与dtm结合，解决了删除缓存未能解决的版本不一致问题
- 发起一个请求，演示了普通删除缓存方案下，版本不一致的问题 `curl http://localhost:8081/api/busi/version?mode=delete`
- 发起一个请求，演示了rockscache+dtm配合解决版本不一致的问题 `curl http://localhost:8081/api/busi/version?mode=rockscache`

在这个demo里面，主要有以下几点
1. 会将DB的数据初始化为v1，然后查询缓存，无数据后查询DB获得v1，睡眠几秒，模拟网络延迟，然后更新缓存。
2. 接着将数据修改为v2，然后查询缓存，无数据后查询DB获得v2，睡眠几毫秒，模拟无网络延迟，然后更新缓存。

- 对于mode=delete，由于网络延迟，导致v1写入缓存的时间在v2之后，导致缓存中的最终版本为v1，与数据库不一致
- 对于mode=rockscache，虽然遇见了网络延时，但是最终缓存中的版本是v2，与数据库中的一致
### 保证原子性
- 发起一个普通更新数据请求，模拟crash，导致DB与缓存不一致 `curl http://localhost:8081/api/busi/atomic?mode=none`
- 发起一个通过dtm更新数据请求，模拟crash，导致DB与缓存不一致，但是5s后，DB与缓存恢复一致 `curl http://localhost:8081/api/busi/atomic?mode=dtm`

可以看到通过dtm更新缓存，在发生进程crash的情况下，也能够保证更新DB和缓存同时成功，或者同时失败。

这个例子也演示了dtm不仅可以和rockscache组合使用，也可以独立于rockscache使用，架构方案会比本地消息表、事务消息、binlog监听的这些架构更加简单

### 强一致访问
rockscache+dtm 也能够提供强一致的访问，访问方式不变，仅需要对rockscache的client设置`StringConsistency`选项

发起一个演示强一致访问的例子 `curl http://localhost:8081/api/busi/strong`，该例子会做以下事情
1. 初始化数据，包括数据库和缓存
2. 采用dtm二阶段消息方式更新数据，保证DB与缓存操作的原子性，其中缓存数据需要3s计算时间
3. 在缓存更新完成前，查询全局事务状态，状态为未完成。（如果用户查询业务结果，即使DB更新完成，只要全局事务未完成，都要告知用户未完成）
4. 查询缓存，在强一致的访问模式下，该查询会等待2中的结果，与最终一致不同


### 升降级时的强一致
如果对于升降级这个短暂的时间窗口，也要求保持强一致，rockscache+dtm也可以做到。

发起一个演示升降级中也保持强一致的例子`curl http://localhost:8081/api/busi/downgrade`，该例子会做以下事情
1. 初始化数据，包括数据库和缓存
2. 采用dtm的Saga模式来更新数据，保证锁缓存、更新DB、删除缓存三个操作的原子性，另外佳鼎缓存数据需要3s计算时间
3. 更新完DB后，还未更新缓存，即可告诉用户业务已完成（与前面强一致的案例不同，条件放松）
4. 查询缓存，在强一致的访问模式下，该查询会等待2中的结果
5. 降级时，先关闭读缓存，等待所有读都不访问缓存时，然后关闭删除缓存
6. 升级时，先打开删除缓存，保证所有的DB更新都会更新到缓存，然后打开读缓存

说明：2中因为要考虑更新DB可能会失败，所以采用SAGA事务模式，而不是消息
