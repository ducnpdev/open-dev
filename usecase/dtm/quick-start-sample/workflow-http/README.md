English | [简体中文](./README-cn.md)

# workflow-http
Minimal usage example for client/workflow using http protocol

## Quick start

### Installing and running dtm

Refer to [dtm installation and running](https://en.dtm.pub/guide/install.html)

### Startup example

``` bash
go run main.go
```

### Output

The order of execution can be seen in the log of workflow-http as follows.

- TransOut
- TransIn

The entire workflow transaction was executed successfully

## Steps
A complete example includes following steps:

### Init Workflow
``` Go
	app.POST(qsBusiAPI+"/workflowResume", func(ctx *gin.Context) {
		log.Printf("workflowResume")
		data, err := ioutil.ReadAll(ctx.Request.Body)
		logger.FatalIfError(err)
		workflow.ExecuteByQS(ctx.Request.URL.Query(), data)
	})

	workflow.InitHTTP(dtmServer, qsBusi+"/workflowResume")
```

### Register a Workflow
``` Go
	wfName := "workflow-http"
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
	// ...
		_, err = wf.NewBranch().NewRequest().SetBody(req).Post(qsBusi + "/TransOut")
	// http request using wf.NewRequest() will be automaticly be intercepted and recorded
	}
```

### Execute a Workflow
``` Go
	err = workflow.Execute(wfName, shortuuid.New(), data)
```

### For more examples, see [dtm-examples](https://github.com/dtm-labs/dtm-examples)
