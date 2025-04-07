package conf

import (
	"fmt"

	"github.com/dtm-labs/dtmcli"
)

var DBConf = dtmcli.DBConf{
	Driver:   "mysql",
	Host:     "en.dtm.pub",
	User:     "dtm",
	Password: "passwd123dtm",
	Port:     3306,
}

var DtmServer = "http://localhost:36789/api/dtmsvr"

const BusiAPI = "/api/busi"
const BusiPort = 8081

var BusiUrl = fmt.Sprintf("http://localhost:%d%s", BusiPort, BusiAPI)
