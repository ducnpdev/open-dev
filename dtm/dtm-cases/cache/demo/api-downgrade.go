package demo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dtm-labs/dtm-cases/utils"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmcli/logger"
	"github.com/dtm-labs/rockscache"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
)

func init() {
	strong := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	strong.Options.StrongConsistency = true

	BusiApp.GET(BusiAPI+"/downgrade", utils.WrapHandler(func(c *gin.Context) interface{} {
		initData(DataKey, "v1", "rockscache")
		data := "v2"
		mode := "rockscache"
		gid := shortuuid.New()
		go Post(BusiUrl+"/downgradeUpdateData", map[string]interface{}{
			"key":       DataKey,
			"value":     data,
			"mode":      mode,
			"time_cost": "3s",
			"gid":       gid,
		})
		time.Sleep(1500 * time.Millisecond)
		r := Get(BusiUrl + "/downgradeQueryProgress?" + fmt.Sprintf("key=%s&value=%s", DataKey, data))
		ensure(r.String() == "\"finished\"", "business status has finished")
		v, err := strong.Fetch(DataKey, 300*time.Second, func() (string, error) {
			return GetDBValue(DataKey).V, nil
		})
		logger.FatalIfError(err)
		ensure(v == "v2", "db value should be v2")

		strong.Options.DisableCacheRead = true
		v, err = strong.Fetch(DataKey, 300*time.Second, func() (string, error) {
			return GetDBValue(DataKey).V, nil
		})
		logger.FatalIfError(err)
		ensure(v == "v2", "db value should be v2")

		logger.Infof("sleeping 3s for all cache read been disabled")
		time.Sleep(3 * time.Second)

		strong.Options.DisableCacheDelete = true
		data = "v3"
		gid2 := shortuuid.New()
		go Post(BusiUrl+"/downgradeUpdateData", map[string]interface{}{
			"key":       DataKey,
			"value":     data,
			"mode":      mode,
			"time_cost": "5ms",
			"gid":       gid2,
		})
		time.Sleep(1500 * time.Millisecond)
		v, err = strong.Fetch(DataKey, 300*time.Second, func() (string, error) {
			return GetDBValue(DataKey).V, nil
		})
		logger.FatalIfError(err)
		ensure(v == data, "db value should be v3")

		logger.Infof("sleeping 3s for all cache delete been disabled")
		time.Sleep(3 * time.Second)
		logger.Infof("now cache can be removed")
		return "ok"
	}))

	BusiApp.GET(BusiAPI+"/downgradeQueryProgress", utils.WrapHandler(func(c *gin.Context) interface{} {
		key := c.Query("key")
		value := c.Query("value")
		evalue := GetDBValue(key).V
		if evalue == value {
			return "finished"
		}
		return "running"
	}))

	BusiApp.POST(BusiAPI+"/downgradeUpdateData", utils.WrapHandler(func(c *gin.Context) interface{} {
		body := MustMapBodyFrom(c)
		saga := dtmcli.NewSaga(DtmServer, body["gid"].(string)).
			Add(BusiUrl+"/downgradeBranch", "", body)
		saga.RetryInterval = 3
		return saga.Submit()
	}))
	BusiApp.POST(BusiAPI+"/downgradeBranch", utils.WrapHandler(func(c *gin.Context) interface{} {
		body := MustMapBodyFrom(c)
		err := strong.LockForUpdate(body["key"].(string), body["gid"].(string))
		if err != nil {
			return err
		}
		bb, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			return err
		}
		err = bb.CallWithDB(db, func(tx *sql.Tx) error {
			// if business failed, user should return error dtmcli.ErrFailure
			// other error will be retried
			return UpdateInTx(tx, &DBRow{
				K:        body["key"].(string),
				V:        body["value"].(string),
				TimeCost: body["time_cost"].(string),
			})
		})
		if err != nil && !errors.Is(err, dtmcli.ErrFailure) {
			return err
		}
		err2 := strong.UnlockForUpdate(body["key"].(string), body["gid"].(string))
		if err2 != nil {
			return err2
		}
		return err
	}))
}
