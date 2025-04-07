package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dtm-labs/dtm-cases/order/common"
	"github.com/dtm-labs/dtm-cases/order/conf"
	"github.com/dtm-labs/dtm-cases/order/service"
	"github.com/dtm-labs/dtm-cases/utils"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmcli/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // register mysql driver
	"github.com/lithammer/shortuuid/v3"
)

func main() {
	logger.InitLog("debug")
	startSvr()
	// fireRequest(defaultReq())
	select {}
}

func startSvr() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	addRoutes(app)
	log.Printf("order examples listening at %d", conf.BusiPort)
	go app.Run(fmt.Sprintf(":%d", conf.BusiPort))
	time.Sleep(100 * time.Millisecond)
}

func addRoutes(app *gin.Engine) {
	service.AddAPIRoute(app)
	service.AddCouponRoute(app)
	service.AddOrderRoute(app)
	service.AddPayRoute(app)
	service.AddStockRoute(app)
	app.Any("/api/fireSucceed", utils.WrapHandler(func(c *gin.Context) interface{} {
		req := defaultReq()
		return fireRequest(req)
	}))
	app.Any("/api/fireFailed", utils.WrapHandler(func(c *gin.Context) interface{} {
		req := defaultReq()
		req.ProductCount = 1000
		return fireRequest(req)
	}))
	app.Any("/api/fireFailedCoupon", utils.WrapHandler(func(c *gin.Context) interface{} {
		req := defaultReq()
		req.CouponID = 101
		return fireRequest(req)
	}))
}

func defaultReq() *common.Req {
	return &common.Req{
		UserID:       1,
		OrderID:      shortuuid.New(),
		ProductID:    1,
		ProductCount: 1,
		CouponID:     0,
		Amount:       100,
	}
}

func fireRequest(req *common.Req) interface{} {
	resty := dtmcli.GetRestyClient()
	resp, err := resty.R().SetBody(req).Post(conf.BusiUrl + "/submitOrder")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return resp.Error()
	}
	return &req
}
