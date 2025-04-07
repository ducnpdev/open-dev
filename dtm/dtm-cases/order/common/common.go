package common

import (
	"database/sql"

	"github.com/dtm-labs/dtm-cases/order/conf"
	"github.com/dtm-labs/dtmcli/dtmimp"
	"github.com/gin-gonic/gin"
)

func DBGet() *sql.DB {
	db, err := dtmimp.PooledDB(conf.DBConf)
	if err != nil {
		panic(err)
	}
	return db
}

type Req struct {
	UserID       int    `json:"user_id"`
	OrderID      string `json:"order_id"`
	Amount       int    `json:"amount"` // money amount
	ProductID    int    `json:"product_id"`
	ProductCount int    `json:"product_count"` // how many product to order
	CouponID     int    `json:"coupon_id"`     // optional
}

func MustGetReq(c *gin.Context) *Req {
	var req Req
	err := c.BindJSON(&req)
	if err != nil {
		panic(err)
	}
	if req.UserID == 0 {
		panic("user_id not specified")
	}
	if req.OrderID == "" {
		panic("order_Id not specified")
	}
	if req.Amount == 0 {
		panic("amount not specified")
	}
	if req.ProductID == 0 {
		panic("product_id not specified")
	}
	return &req
}
