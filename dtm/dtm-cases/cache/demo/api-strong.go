package demo

import (
	"database/sql"
	"encoding/json"
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

	BusiApp.GET(BusiAPI+"/strong", utils.WrapHandler(func(c *gin.Context) interface{} {
		initData(DataKey, "v1", "rockscache")
		data := "v2"
		mode := "rockscache"
		gid := shortuuid.New()
		Post(BusiUrl+"/strongUpdateData", map[string]interface{}{
			"key":       DataKey,
			"value":     data,
			"mode":      mode,
			"time_cost": "3s",
			"gid":       gid,
		})
		r := Get(BusiUrl + "/strongQueryProgress?gid=" + gid)
		ensure(r.String() == "\"submitted\"", "business status has not finished")
		v, err := strong.Fetch(DataKey, 300*time.Second, func() (string, error) {
			return GetDBValue(DataKey).V, nil
		})
		logger.FatalIfError(err)
		ensure(v == "v2", "db value should be v2")
		return "ok"
	}))

	BusiApp.GET(BusiAPI+"/strongQueryProgress", utils.WrapHandler(func(c *gin.Context) interface{} {
		gid := c.Query("gid")
		r := Get(DtmServer + "/query?gid=" + gid)
		var m map[string]interface{}
		err := json.Unmarshal(r.Body(), &m)
		logger.FatalIfError(err)
		return m["transaction"].(map[string]interface{})["status"].(string)
	}))

	BusiApp.POST(BusiAPI+"/strongUpdateData", utils.WrapHandler(func(c *gin.Context) interface{} {
		body := MustMapBodyFrom(c)
		gid := body["gid"].(string)
		msg := dtmcli.NewMsg(DtmServer, gid).
			Add(BusiUrl+"/deleteCache", body)
		msg.WaitResult = false // when return success, the global transaction has finished
		return msg.DoAndSubmit(BusiUrl+"/queryPrepared", func(bb *dtmcli.BranchBarrier) error {
			return bb.CallWithDB(db, func(tx *sql.Tx) error {
				return UpdateInTx(tx, &DBRow{
					K:        body["key"].(string),
					V:        body["value"].(string),
					TimeCost: body["time_cost"].(string),
				})
			})
		})
	}))
}
