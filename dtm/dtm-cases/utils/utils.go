package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmcli/logger"
	"github.com/gin-gonic/gin"
)

// WrapHandler used by examples. much more simpler than WrapHandler2
func WrapHandler(fn func(*gin.Context) interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		began := time.Now()
		ret := fn(c)
		status, res := dtmcli.Result2HttpJSON(ret)

		b, _ := json.Marshal(res)
		if status == http.StatusOK || status == http.StatusTooEarly {
			logger.Infof("%2dms %d %s %s %s", time.Since(began).Milliseconds(), status, c.Request.Method, c.Request.RequestURI, string(b))
		} else {
			logger.Errorf("%2dms %d %s %s %s", time.Since(began).Milliseconds(), status, c.Request.Method, c.Request.RequestURI, string(b))
		}
		c.JSON(status, res)
	}
}

func MustBarrierFrom(c *gin.Context) *dtmcli.BranchBarrier {
	bb, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
	if err != nil {
		panic(err)
	}
	return bb
}
