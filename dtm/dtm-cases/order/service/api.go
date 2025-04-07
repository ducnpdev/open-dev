package service

import (
	"github.com/dtm-labs/dtm-cases/order/common"
	"github.com/dtm-labs/dtm-cases/order/conf"
	"github.com/dtm-labs/dtm-cases/utils"
	"github.com/dtm-labs/dtmcli"
	"github.com/gin-gonic/gin"
)

func AddAPIRoute(app *gin.Engine) {
	app.POST("/api/busi/submitOrder", utils.WrapHandler(func(c *gin.Context) interface{} {
		req := common.MustGetReq(c)
		saga := dtmcli.NewSaga(conf.DtmServer, "gid-"+req.OrderID).
			Add(conf.BusiUrl+"/orderCreate", conf.BusiUrl+"/orderCreateRevert", &req).
			Add(conf.BusiUrl+"/stockDeduct", conf.BusiUrl+"/stockDeductRevert", &req).
			Add(conf.BusiUrl+"/couponUse", conf.BusiUrl+"couponUseRevert", &req).
			Add(conf.BusiUrl+"/payCreate", conf.BusiUrl+"/payCreateRevert", &req)
		return saga.Submit()
	}))
}
