package service

import (
	"database/sql"

	"github.com/dtm-labs/dtm-cases/order/common"
	"github.com/dtm-labs/dtm-cases/utils"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmcli/dtmimp"
	"github.com/gin-gonic/gin"
)

func AddStockRoute(app *gin.Engine) {
	app.POST("/api/busi/stockDeduct", utils.WrapHandler(func(c *gin.Context) interface{} {
		req := common.MustGetReq(c)
		return utils.MustBarrierFrom(c).CallWithDB(common.DBGet(), func(tx *sql.Tx) error {
			affected, err := dtmimp.DBExec(tx,
				"update ord.stock set stock=stock-?, update_time=now() where product_id=? and stock >= ?",
				req.ProductCount, req.ProductID, req.ProductCount)
			if err == nil && affected == 0 {
				return dtmcli.ErrFailure // not enough stock, return Failure to rollback
			}
			return err
		})
	}))
	app.POST("/api/busi/stockDeductRevert", utils.WrapHandler(func(c *gin.Context) interface{} {
		req := common.MustGetReq(c)
		return utils.MustBarrierFrom(c).CallWithDB(common.DBGet(), func(tx *sql.Tx) error {
			_, err := dtmimp.DBExec(tx,
				"update ord.stock set stock=stock+?, update_time=now() where product_id=?",
				req.ProductCount, req.ProductID)
			return err
		})
	}))
}
