package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/logger"
	"github.com/dtm-labs/client/workflow"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
)

const qsBusiAPI = "/api/busi_start"
const qsBusiPort = 8082

var qsBusi = fmt.Sprintf("http://localhost:%d%s", qsBusiPort, qsBusiAPI)

func main() {
	var x chan bool
	QsStartSvr()
	_ = QsFireRequest()
	<-x
}

// QsStartSvr quick start: start server
func QsStartSvr() {
	app := gin.New()
	qsAddRoute(app)
	log.Printf("quick start examples listening at %d", qsBusiPort)
	go func() {
		_ = app.Run(fmt.Sprintf(":%d", qsBusiPort))
	}()
}

func qsAddRoute(app *gin.Engine) {
	app.POST(qsBusiAPI+"/TransIn", func(c *gin.Context) {
		log.Printf("TransIn")
		time.Sleep(time.Second)
		c.JSON(200, "")
		// c.JSON(409, "") // Status 409 for Failure. Won't be retried
	})
	app.POST(qsBusiAPI+"/TransInCompensate", func(c *gin.Context) {
		log.Printf("TransInCompensate")
		time.Sleep(time.Second)
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/TransOut", func(c *gin.Context) {
		// TODO get value
		// var body map[string]interface{}
		// if err := c.ShouldBindJSON(&body); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// amount, ok := body["amount"].(float64) // JSON numbers are parsed as float64 by default
		// if !ok {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount type"})
		// 	return
		// }
		// fmt.Println("amount:", amount)

		//
		log.Printf("TransOut")
		time.Sleep(time.Second)
		c.JSON(200, "")
		// c.JSON(409, "") // Status 409 for Failure. Won't be retried
	})
	app.POST(qsBusiAPI+"/TransOutCompensate", func(c *gin.Context) {
		log.Printf("TransOutCompensate")
		time.Sleep(time.Second)
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/notify", func(c *gin.Context) {
		log.Printf("notify")
		time.Sleep(time.Second)
		// c.JSON(200, "")
		c.JSON(409, "") // Status 409 for Failure. Won't be retried
	})
	app.POST(qsBusiAPI+"/notifyCompensate", func(c *gin.Context) {
		log.Printf("notifyCompensate")
		time.Sleep(time.Second)
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/workflowResume", func(ctx *gin.Context) {
		log.Printf("workflowResume")
		data, err := ioutil.ReadAll(ctx.Request.Body)
		logger.FatalIfError(err)
		workflow.ExecuteByQS(ctx.Request.URL.Query(), data)
	})
}

const dtmServer = "http://localhost:36789/api/dtmsvr"

// QsFireRequest quick start: fire request
func QsFireRequest() string {
	workflow.InitHTTP(dtmServer, qsBusi+"/workflowResume")
	wfName := "workflow-http"
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var req gin.H
		err := json.Unmarshal(data, &req)
		logger.FatalIfError(err)
		wf.BranchID = "01"
		_, err = wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, err := wf.NewRequest().SetBody(req).Post(qsBusi + "/TransOutCompensate")
			return err
		}).NewRequest().SetBody(req).Post(qsBusi + "/TransOut")
		if err != nil {
			return err
		}
		_, err = wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, err := wf.NewRequest().SetBody(req).Post(qsBusi + "/TransInCompensate")
			return err
		}).NewRequest().SetBody(req).Post(qsBusi + "/TransIn")
		if err != nil {
			return err
		}
		_, err = wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, err := wf.NewRequest().SetBody(req).Post(qsBusi + "/notifyCompensate")
			return err
		}).NewRequest().SetBody(req).Post(qsBusi + "/notify")
		return err
	})
	logger.FatalIfError(err)

	gid := shortuuid.New()
	req := &gin.H{"amount": 30} // the payload of requests
	data, err := json.Marshal(req)
	logger.FatalIfError(err)
	err = workflow.Execute(wfName, gid, data)
	logger.Infof("workflow.Execute %s result is: %v", gid, err)
	return gid
}
