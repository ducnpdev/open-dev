# workflow-http
client/workflow的最简http使用示例

## 快速开始

### 安装运行dtm

参考[dtm安装运行](https://dtm.pub/guide/install.html)

### 启动示例

``` bash
go run main.go
```

### 输出

可以从workflow-http的日志里看到执行的顺序如下：

- TransOut
- TransIn

整个workflow事务执行成功

## 步骤
完整的使用例子会包括以下步骤

### 初始化 Workflow
``` Go
	app.POST(qsBusiAPI+"/workflowResume", func(ctx *gin.Context) {
		log.Printf("workflowResume")
		data, err := ioutil.ReadAll(ctx.Request.Body)
		logger.FatalIfError(err)
		workflow.ExecuteByQS(ctx.Request.URL.Query(), data)
	})

	workflow.InitHTTP(dtmServer, qsBusi+"/workflowResume")
```

### 注册一个workflow
``` Go
	wfName := "workflow-http"
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
	// ...
		_, err = wf.NewBranch().NewRequest().SetBody(req).Post(qsBusi + "/TransOut")
	// 使用 wf.NewRequest() 的http请求，会被自动拦截，并且记录进度
	}
```

### 执行一个 Workflow
``` Go
	err = workflow.Execute(wfName, shortuuid.New(), data)
```

### 更多示例，详见[dtm-examples](https://github.com/dtm-labs/dtm-examples)
