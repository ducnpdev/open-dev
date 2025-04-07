English | [简体中文](./README-cn.md)

# Order applications

This example can be read in conjunction with the [Order Application](https://en.dtm.pub/app/order.html) in the dtm documentation

## Overview
This project demonstrates how dtm can be applied to a non-monolithic order system to ensure that multiple steps in an order are eventually executed "atomically" to ensure eventual consistency.

#### Starting dtm
[Quick start dtm](https://en.dtm.pub/guide/install.html)

#### Run this example
`go run main.go`

#### Initiating an order request
- Initiate a normal order `curl http://localhost:8081/api/fireSucceed`
- Launch a rollback order due to insufficient stock `curl http://localhost:8081/api/fireFailed`
- Initiating a rollback order for a failed coupon deduction `curl http://localhost:8081/api/fireFailedCoupon`

The following points should be noted.
1. fireFailed requests, which fail because of insufficient stock, will be rolled back by the global transaction. The rollback is a null compensation and no stock adjustment will be done.
2. fireFailedCoupon requests, which fail because of unsuccessful coupon deductions, are rolled back. The stock will also be rolled back and real stock adjustment will be done.
3. developers do not need to care about null compensation, developers only need to care about how to deduct inventory and roll back inventory, whether null compensation, and how to roll back null compensation, etc., will be automatically handled by the dtm framework.

#### Directory Description
This project has the following contents
- main.go: the main program file
- service directory: the relevant individual service files
- conf directory: the relevant configuration
- common directory: code shared by multiple services
- order.sql: the order system table needed to create this example
