简体中文 | [English](./README.md)

# 订单应用

此项目可以结合dtm文档中的 [订单应用](https://dtm.pub/app/order.html)阅读

## 概述
本项目主要演示了dtm如何应用于非单体的订单系统，保证订单中的多个步骤，能够最终“原子”执行，保证最终一致性。

#### 启动dtm
[快速启动dtm](https://dtm.pub/guide/install.html)

#### 运行本例子
`go run main.go`

#### 发起订单请求
- 发起一个正常订单 `curl http://localhost:8081/api/fireSucceed`
- 发起一个因库存不足的回滚订单 `curl http://localhost:8081/api/fireFailed`
- 发起一个因扣减优惠券失败而回滚的订单 `curl http://localhost:8081/api/fireFailedCoupon`

以下几点说明一下：
1. fireFailed 请求，是因为库存不足而失败，此时全局事务会回滚。回滚时会进行库存的回滚操作，此时库存回滚时发生了一个空补偿，在实际操作中不会进行库存相关的业务操作
2. fireFailedCoupon 请求，是因为扣减优惠券不成功而失败，回滚时会进行库存的回滚操作，此时库存回滚时发生了一个正常补偿，在实际操作中会进行库存相关的业务操作
3. 开发人员无需关心是否空补偿，开发者只需要关心如何扣减库存和回滚库存，是否空补偿，以及如何回滚空补偿等，都会有dtm框架进行自动处理。

#### 目录说明
本项目有以下内容
- main.go: 主程序文件
- service目录：相关各个服务文件
- conf目录：相关配置
- common目录：多个服务共享的代码
- order.sql: 创建本例子所需的订单系统表
