English | [简体中文](./README-cn.md)

# Flash-sale application
This project can be read in conjunction with the [flash-sale](https://en.dtm.pub/app/flash.html) in the dtm documentation

## Overview
This project demonstrates how dtm can be applied to a flash-sale system, and will demonstrate how a flash-sale system can ensure accurate stock deductions and create accurate quantities of orders even when a process crash occurs.

#### start dtm
[Quick start dtm](https://en.dtm.pub/guide/install.html)

#### Run this example
`go run main.go`

#### Launching an order request
- Launch a flash-sale request that completes normally `curl http://localhost:8081/api/busi/flashSales`
- Launch a flash-sale request that crashes when the stock deduction is complete `curl http://localhost:8081/api/busi/flashSales-crash` Wait about ten seconds or so for the order to be created without impact

#### Simulates a flash-sale
`curl http://localhost:8081/api/busi/flashSales-batch`

After the user initiates this simulated flash-sale request, the example does the following.
1. Reset the variables: set the stock to 4 and reset the number of orders created to 0
2. Initiate a flash-sale request. The handler of this request will simulate a crash after the stock is deducted. It will then sleep for 0.5s to ensure that the stock is deducted to 3 before proceeding to the next step
3. Launch 10 concurrent flash-sale requests, then sleep for 0.5s to ensure that the 10 flash-sale requests are processed before proceeding to the next step
4. Output the number of orders created at this point as 3, all three of which were generated in step 3
5. Wait about 3-5 seconds and you can see that the number of orders created changes to 4. This is due to the global transaction timeout check in step 2, which finally succeeds and creates the 4th order

Conclusion: The approach in this example ensures that the inventory in Redis and the orders in the database end up being strictly consistent, eliminating the need to manually calibrate the data
