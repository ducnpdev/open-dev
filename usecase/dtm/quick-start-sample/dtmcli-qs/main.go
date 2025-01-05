package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
)

const qsBusiAPI = "/api/busi_start"
const qsBusiPort = 8082

var qsBusi = fmt.Sprintf("http://localhost:%d%s", qsBusiPort, qsBusiAPI)

func main() {
	QsStartSvr()
	_ = QsFireRequest()
	time.Sleep(3 * time.Second)
}

// QsStartSvr quick start: start server
func QsStartSvr() {
	app := gin.New()
	qsAddRoute(app)
	log.Printf("quick start examples listening at %d", qsBusiPort)
	go func() {
		_ = app.Run(fmt.Sprintf(":%d", qsBusiPort))
	}()
	time.Sleep(100 * time.Millisecond)
}

func qsAddRoute(app *gin.Engine) {
	app.POST(qsBusiAPI+"/TransIn", func(c *gin.Context) {
		log.Printf("TransIn")
		c.JSON(200, "")
		// c.JSON(409, "") // Status 409 for Failure. Won't be retried
	})
	app.POST(qsBusiAPI+"/TransInCompensate", func(c *gin.Context) {
		log.Printf("TransInCompensate")
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/TransOut", func(c *gin.Context) {
		log.Printf("TransOut")
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/TransOutCompensate", func(c *gin.Context) {
		log.Printf("TransOutCompensate")
		c.JSON(200, "")
	})
}

const dtmServer = "http://localhost:36789/api/dtmsvr"

// QsFireRequest quick start: fire request
func QsFireRequest() string {
	req := &gin.H{"amount": 30} // the payload of requests
	gid := shortuuid.New()
	// DtmServer is the address of dtm
	saga := dtmcli.NewSaga(dtmServer, gid).
		// add a branch transaction，action url is: qsBusi+"/TransOut"， compensate url: qsBusi+"/TransOutCompensate"
		Add(qsBusi+"/TransOut", qsBusi+"/TransOutCompensate", req).
		// add a branch transaction，action url is: qsBusi+"/TransIn"， compensate url: qsBusi+"/TransInCompensate"
		Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// submit saga global transaction，dtm will finish all action and compensation
	err := saga.Submit()

	if err != nil {
		panic(err)
	}
	log.Printf("transaction: %s submitted", saga.Gid)
	return saga.Gid
}
