package service

import (
	"database/sql"

	"github.com/dtm-labs/dtm-cases/order/common"
	"github.com/dtm-labs/dtm-cases/utils"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmcli/dtmimp"
	"github.com/gin-gonic/gin"
)

func AddCouponRoute(app *gin.Engine) {
	app.POST("/api/busi/couponUse", utils.WrapHandler(func(c *gin.Context) interface{} {
		req := common.MustGetReq(c)
		if req.CouponID == 0 { // no coupon to use, so return
			return nil
		}
		bb := utils.MustBarrierFrom(c)
		return bb.CallWithDB(common.DBGet(), func(tx *sql.Tx) error {
			affected, err := dtmimp.DBExec(tx,
				"update ord.user_coupon set used=1, update_time=now() where used=0 and use_id=? and coupon_id=?",
				req.UserID, req.CouponID)
			if err == nil && affected == 0 {
				return dtmcli.ErrFailure // update coupon use failed, return failure to rollback
			}
			return err
		})
	}))
	app.POST("/api/busi/couponUseRevert", utils.WrapHandler(func(c *gin.Context) interface{} {
		req := common.MustGetReq(c)
		if req.CouponID == 0 { // no coupon to use, so return
			return nil
		}
		bb := utils.MustBarrierFrom(c)
		return bb.CallWithDB(common.DBGet(), func(tx *sql.Tx) error {
			_, err := dtmimp.DBExec(tx,
				"update ord.user_coupon set used=0, update_time=now() where used=1 and use_id=? and coupon_id=?",
				req.UserID, req.CouponID)
			return err
		})
	}))
}
