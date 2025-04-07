package demo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dtm-labs/dtm-cases/utils"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmcli/logger"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
)

func init() {
	BusiApp.GET(BusiAPI+"/atomic", utils.WrapHandler(func(c *gin.Context) interface{} {
		mode := c.Query("mode")
		if mode != "normal" && mode != "dtm" {
			return errors.New("mode should be normal or dtm")
		}
		initData(DataKey, "v1", mode)

		data := "v2"
		_ = Post(BusiUrl+"/atomicCrashUpdate", map[string]interface{}{
			"key":       DataKey,
			"value":     data,
			"mode":      mode,
			"time_cost": "5ms",
		})
		logger.Infof("sleeping for 7 seconds to wait for the update")
		time.Sleep(7 * time.Second)

		_, _ = Fetch(mode, DataKey, 300*time.Millisecond, func() (string, error) {
			return GetDBValue(DataKey).V, nil
		})
		dbv := GetDBValue(DataKey)
		cachev := GetCacheValue(DataKey, mode)
		ensure(dbv.V == "v2", "db value should be v2")
		if mode == "normal" {
			ensure(cachev == "v1", "cache value should be v1")
			logger.Infof("for mode normal, db value is: %s, cache value is: %s, not matched", dbv, cachev)
		} else {
			ensure(cachev == "v2", "cache value should be v2")
			logger.Infof("for mode dtm, db value is: %s, cache value is: %s, matched", dbv, cachev)
		}
		return "ok"
	}))

	BusiApp.POST(BusiAPI+"/atomicCrashUpdate", utils.WrapHandler(func(c *gin.Context) interface{} {
		body := MustMapBodyFrom(c)
		mode := body["mode"].(string)
		if mode != "normal" && mode != "dtm" {
			return errors.New("mode should be normal or dtm")
		}
		row := &DBRow{
			K:        body["key"].(string),
			V:        body["value"].(string),
			TimeCost: body["time_cost"].(string),
		}
		go func() {
			if mode == "normal" {
				SetDBValue(row)
				time.Sleep(86400 * time.Second) // simulate a crash
				DeleteCacheValue(body["key"].(string))
			} else {
				msg := dtmcli.NewMsg(DtmServer, shortuuid.New()).
					Add(BusiUrl+"/deleteCache", map[string]interface{}{
						"key":  body["key"].(string),
						"mode": mode,
					})
				msg.TimeoutToFail = 3
				err := msg.DoAndSubmit(BusiUrl+"/queryPrepared", func(bb *dtmcli.BranchBarrier) error {
					err := bb.CallWithDB(db, func(tx *sql.Tx) error {
						return UpdateInTx(tx, row)
					})
					time.Sleep(86400 * time.Second) // simulate a crash
					return err
				})
				logger.FatalIfError(err)
			}
		}()
		time.Sleep(200 * time.Millisecond) // wait for db update
		return "ok"
	}))

}
