package sqs

import "fmt"

func MainSQS() {
	queueName := "sqs-demo"
	ses := GetSession()
	list, err := ListQueues(ses)
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
	queueUrl, err := GetQueueURL(ses, queueName)
	if err != nil {
		panic(err)
	}
	fmt.Println(*queueUrl.QueueUrl)

	// _, err = SendMessage(ses, *queueUrl.QueueUrl, "this is test message")
	// if err != nil {
	// 	panic(err)
	// }

	listM, err := GetMessages(ses, *queueUrl.QueueUrl, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(listM)
	for _, item := range listM.Messages {
		err = DeleteMessage(ses, *queueUrl.QueueUrl, item.ReceiptHandle)
		if err != nil {
			panic(err)
		}
	}

}
